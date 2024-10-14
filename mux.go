package main

import (
	"net/http"

	"github.com/YutaKakiki/go-todo-api/handler"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux() http.Handler {
	// "/health" エンドポイントに対するハンドラを設定
	mux := chi.NewRouter()
	// *chi.Muxはhandler型を満たすので、そのまま使える！
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json;charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)
	return mux
}
