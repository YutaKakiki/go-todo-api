package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/go-playground/validator/v10"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

// Handler型を満たす
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストスコープにおいてcontextを生成
	ctx := r.Context()
	// リクエストボディの構造
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	// jsonをデコード（受け取る）して変数bに格納
	// json⇨構造体
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		// デコード時にエラー出たら
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 構造体のバリデーション
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// ボディにあるタイトルから、新たなentity.Taskオブジェクトを生成
	t := &entity.Task{
		Title:   b.Title,
		Status:  entity.TaskStatusTodo, //定数でステータスは定義してある
		Created: time.Now(),
	}
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: id}
	RespondJSON(ctx, w, rsp, http.StatusOK)

}
