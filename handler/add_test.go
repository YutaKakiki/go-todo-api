package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YutaKakiki/go-todo-api/entity"
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
			// モックを使って実際のDBには干渉しない
			moq := &AddTaskServiceMock{}

			// AddTaskメソッドの挙動を記述
			// ビジネスロジックを担当するserviceは通常DB操作を行うが、
			// テストで実際のDBと干渉したくないため、モックを使ってDB操作の挙動を模倣する
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				// HTTPリクエストの期待するステータスコードが200（成功）である場合
				// モックされたAddTaskメソッドは、ID: 1 のタスクを返す
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				// それ以外の場合は、エラーメッセージを返す
				return nil, errors.New("error from mock")
			}

			// SUT (System Under Test) は、現在テストされているシステム
			// 今回のテスト対象はAddTaskハンドラなので、それをセットアップする
			sut := AddTask{
				Service:   moq,             // サービスにはモックを注入
				Validator: validator.New(), // バリデータをセットアップ
			}

			// ServeHTTPメソッドを実行して、ハンドラの動作をテスト
			// この中でモックされたAddTaskメソッドが呼び出される
			sut.ServeHTTP(w, r)

			// レスポンスを受け取る
			rsp := w.Result()
			// 期待するレスポンスと実際のレスポンスを比較
			testutil.AsertResponse(t, rsp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}

}
