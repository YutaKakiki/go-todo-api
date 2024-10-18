package store

import (
	"context"
	"errors"
	"testing"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/testutil"
)

func TestKVS_Save(t *testing.T) {
	cli := testutil.OpenRedisForTest(t)

	sut := &KVS{Cli: cli}
	ctx := context.Background()
	key := "TestKey"
	uid := entity.UserID(1234)
	// キーを削除しておく（あったら）
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	err := sut.Save(ctx, key, uid)
	if err != nil {
		t.Errorf("want no error,but got :%v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{Cli: cli}
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		key := "TestKey_OK"
		uid := entity.UserID(1234)
		// キーを削除しておく（あったら）
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		sut.Save(ctx, key, uid)
		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want to error,but got %v", err)
		}
		if got != uid {
			t.Errorf("want %d,but got %d", uid, got)
		}
	})

	t.Run("not Found", func(t *testing.T) {
		t.Parallel()

		key := "TestKey_NotFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v,but got %v(value = %d)", ErrNotFound, err, got)
		}

	})
}
