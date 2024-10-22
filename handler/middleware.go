package handler

import (
	"net/http"

	"github.com/YutaKakiki/go-todo-api/auth"
)

// jwから情報を埋め込むミドルウェア
func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// リクエストヘッダからjwtを取得し、そこから得た情報をリクエストスコープのコンテキストに付加
			// 内部で検証も行われている
			req, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

// 上のミドルウェアを通過した後（contextに情報が付加された後）
// adminなのか検証する

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
