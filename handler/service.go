package handler

import (
	"context"

	"github.com/YutaKakiki/go-todo-api/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService RegisterUserService LoginService
type ListTasksService interface {
	ListTask(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}
