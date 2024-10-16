package service

import (
	"context"
	"fmt"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTask(ctx context.Context) (entity.Tasks, error) {
	ts, err := l.Repo.ListTask(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list :%w", err)
	}
	return ts, nil
}
