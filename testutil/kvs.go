package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()

	host := "127.0.0.1"
	port := 36379
	if _, defined := os.LookupEnv("CI"); defined {
		port = 6379
	}

	// redisクライアントを作成
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})
	// 接続確認
	// エラー出てたら吐きたす
	// 接続されてたらPONGがかえってくるらしい
	if err := cli.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("faile to connect redis:%s", err)
	}
	return cli

}
