package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// http.Server型をラップした構造体定義
type Server struct {
	srv *http.Server
	l   net.Listener
}

// ルーティングの設定を引数で受け取り、ルーティングの責務を取り除く
func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// シグナルを受信するためのコンテキストを作成
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// シグナルの監視をやめてリソースを開放
	defer stop()
	// 引数で受け取ったcontextから新たにキャンセル機能を持つコンテキストを作成
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTP鯖を起動
	// エラーを返す
	eg.Go(func() error {
		// サーバーが正常修了にShutdownされた時を除く
		// Serve:リクエストを処理、レスポンスを返す
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			// エラーは errgroup にわたる
			return err
		}
		return nil
	})

	// 引数として受け取ったctxからキャンセル通知を待つ
	<-ctx.Done()
	// サーバーをシャットダウン
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shutdown: %+v", err)
	}
	// ゴルーチン終了を待つ
	// ゴルーチンがerrを返した時に終了
	// すべてのゴルーチンが正常に終了した場合は nil を返して終了
	return eg.Wait()
}
