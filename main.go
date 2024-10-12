package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/YutaKakiki/go-todo-api/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	// run関数内でサーバーを起動
	if err := run(context.Background()); err != nil {
		fmt.Printf("failed to terminate server:%v", err)
		os.Exit(1)
	}
}

// 外部からのキャンセル操作を受け取るとサーバを修了する
func run(ctx context.Context) error {
	// 環境変数の構造体を取得
	cfg, err := config.New()
	if err != nil {
		return err
	}
	// 環境変数からPORT番号を指定してリッスン
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %v", url)
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
