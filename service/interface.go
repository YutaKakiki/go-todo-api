package service

import (
	"context"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
)

// store pkgに直接依存しないためのインターフェース

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTask(ctx context.Context, db store.Queryer, id entity.UserID) (entity.Tasks, error)
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}
type TokenGenerator interface {
	GenerateToken(ctx context.Context, u *entity.User) ([]byte, error)
}
