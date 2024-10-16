package handler

import (
	"net/http"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/jmoiron/sqlx"
)

type ListTask struct {
	// Store *store.TaskStore
	DB   *sqlx.DB
	Repo *store.Repository
}

// 返ってくるtaskの構造
type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Storeの中のtaskを取得
	// tasks := lt.Store.GetAll()
	tasks, err := lt.Repo.ListTask(ctx, lt.DB)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}
	// 空の構造体を定義（ここに入れてく）
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	// JSON（配列）をレスポンスに書き込む
	RespondJSON(ctx, w, rsp, http.StatusOK)

}
