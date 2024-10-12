package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	// ０：空きポートを割り当て
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port:%v", err)
	}
	// キャンセル機能を持つコンテキストを定義
	ctx, cancel := context.WithCancel(context.Background())
	// run関数をゴルーチンで回してサーバを起動
	// また、返り値のエラーを参照したい
	eg, ctx := errgroup.WithContext(ctx)
	// HTTPリクエストマルチプレクサ:HTTPリクエストを適切なハンドラーにルーティングする
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello,%s", r.URL.Path[1:])
	})
	// HTTP鯖を別ゴルーチンで起動
	eg.Go(func() error {
		s := NewServer(l, mux)
		return s.Run(ctx)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try req to:%q", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// 期待する値と出力を検証
	want := fmt.Sprintf("hello,%s", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	// run関数にキャンセル通知
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

}
