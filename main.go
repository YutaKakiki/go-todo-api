package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/YutaKakiki/go-todo-api/config"
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
	cfg, err := config.New()
	if err != nil {
		return err
	}
	// 環境変数からPORT番号を取得しリッスン
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %v", url)
	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}
