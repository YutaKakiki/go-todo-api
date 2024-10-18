package store

import (
	"context"
	"fmt"
	"time"

	"github.com/YutaKakiki/go-todo-api/config"
	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/redis/go-redis/v9"
)

type KVS struct {
	Cli *redis.Client
}

// Redisへ接続する*store.KVS型のコンストラクタ
func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	// redisクライアントを作成
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
	// 接続確認
	// エラー出てたら吐きたす
	// 接続されてたらPONGがかえってくるらしい
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{Cli: cli}, nil
}

// key:valueを保存する
func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	id := int64(userID)
	return k.Cli.Set(ctx, key, id, 30*time.Minute).Err()
}

// キーを使ってロードする
func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Cli.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), nil
}
