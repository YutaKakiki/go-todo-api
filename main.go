package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	// コマンドの引数からポート番号取得
	p := os.Args[1]
	// リッスン：接続を待ち受ける
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s:%v", p, err)
	}
	// run関数内でサーバーを起動
	if err := run(context.Background(), l); err != nil {
		fmt.Printf("failed to terminate server:%v", err)
	}
}

// 外部からのキャンセル操作を受け取るとサーバを修了する
func run(ctx context.Context, l net.Listener) error {
	s := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello,%s", r.URL.Path[1:])
		}),
	}
	// 引数で受け取ったcontextから新たにキャンセル機能を持つコンテキストを作成
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTP鯖を起動
	// エラーを返す
	eg.Go(func() error {
		// サーバーが正常修了にShutdownされた時を除く
		// Serve:リクエストを処理、レスポンスを返す
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			// エラーは errgroup にわたる
			return err
		}
		return nil
	})

	// 引数として受け取ったctxからキャンセル通知を待つ
	<-ctx.Done()
	// サーバーをシャットダウン
	if err := s.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown: %+v", err)
	}
	// ゴルーチン終了を待つ
	// ゴルーチンがerrを返した時に終了
	// すべてのゴルーチンが正常に終了した場合は nil を返して終了
	return eg.Wait()
}
