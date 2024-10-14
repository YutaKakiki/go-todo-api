package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrResponse struct {
	Message string `json:"message"`
	// フィールドの値が空ならば省略
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	// body（構造体）をJSON形式に変換
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// エラーレスポンスを用意
		rsp := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		// エラーメッセージをレスポンスとしてエンコードし、レスポンスに書き込む
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error response error:%v", err)
		}
		return
	}
	w.WriteHeader(status)
	// レスポンスに書き込む
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
