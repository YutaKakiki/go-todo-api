package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// HTTPサーバーを起動
	err := http.ListenAndServe(
		":8080",
		//Handler型を満たす
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 最初の「/」以降を出力
			fmt.Fprintf(w, "Hello,%s", r.URL.Path[1:])
		}),
	)
	if err != nil {
		fmt.Printf("failed to terminate server:%v", err)
		//異常終了
		os.Exit(1)
	}

}
