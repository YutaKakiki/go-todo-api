package main

import (
	"context"
	"net/http"

	"github.com/YutaKakiki/go-todo-api/clock"
	"github.com/YutaKakiki/go-todo-api/config"
	"github.com/YutaKakiki/go-todo-api/handler"
	"github.com/YutaKakiki/go-todo-api/service"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	// "/health" エンドポイントに対するハンドラを設定
	mux := chi.NewRouter()
	// *chi.Muxはhandler型を満たすので、そのまま使える！
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json;charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clock.RealClocker{}}
	// 埋め込み型によるDI
	at := &handler.AddTask{
		// ServiceフィールドはAddTaskService型：AddTaskメソッドを実装していること
		Service:   &service.AddTask{DB: db, Repo: &r}, //TaskAdder型であるRepoを通して、DBとやり取り
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)
	// 埋め込み型によるDI
	lt := &handler.ListTask{
		Service: &service.ListTask{
			DB:   db,
			Repo: &r,
		},
	}
	mux.Get("/tasks", lt.ServeHTTP)
	ru := &handler.RegisterUser{
		Service: &service.RegisterUser{
			DB:   db,
			Repo: &r,
		},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
