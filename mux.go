package main

import "net/http"

func NewMux() http.Handler {
	// "/health" エンドポイントに対するハンドラを設定
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json;charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
