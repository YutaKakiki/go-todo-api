package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/YutaKakiki/go-todo-api/testutil"
	"github.com/go-playground/validator/v10"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}
	// テーブルドリブン
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"bad": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}

	// テーブルデータを回す：サブテストを並列実行
	for n, tt := range tests {
		// テストデータをバインド
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			// レスポンス、リクエストのモック
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet,
				"/tasks",
				// LoadFileでかえってきたバイトスライスをio.Readerに
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)
			// SUT (System Under Test) は、現在テストされているシステム
			// AddTaskハンドラオブジェクトを設定
			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			// serveHTTPを実行
			sut.ServeHTTP(w, r)
			// レスポンスを受け取る
			rsp := w.Result()
			// 期待するレスポンスと実際のレスポンスを比較
			testutil.AsertResponse(t, rsp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}

}
