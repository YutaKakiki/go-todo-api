package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/store"
	"github.com/YutaKakiki/go-todo-api/testutil"
)

func TestListTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		tasks map[entity.TaskID]*entity.Task
		want  want
	}{
		"ok": {
			tasks: map[entity.TaskID]*entity.Task{
				1: {
					ID:     1,
					Title:  "test1",
					Status: "todo",
				},
				2: {
					ID:     2,
					Title:  "test2",
					Status: "done",
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/ok_rsp.json.golden",
			},
		},
		"empty": {
			tasks: map[entity.TaskID]*entity.Task{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/empty_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			// レスポンス、リクエストのモック
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet,
				"/tasks",
				// LoadFileでかえってきたバイトスライスをio.Readerに
				nil,
			)
			sut := ListTask{
				Store: &store.TaskStore{Tasks: tt.tasks},
			}
			sut.ServeHTTP(w, r)
			rsp := w.Result()
			testutil.AsertResponse(t, rsp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
