package service

import (
	"context"
	"fmt"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

// ハンドラによって呼び出される
func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	// 実際にDBに挿入する
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register :%w", err)
	}
	return t, nil
}
