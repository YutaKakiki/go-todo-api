package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/testutil"
)

func TestListTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		tasks []*entity.Task
		want  want
	}{
		"ok": {
			tasks: []*entity.Task{
				{
					ID:     1,
					Title:  "test1",
					Status: "todo",
				},
				{
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
			tasks: []*entity.Task{},
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
			// モックを用意
			moq := &ListTasksServiceMock{}
			moq.ListTaskFunc = func(ctx context.Context) (entity.Tasks, error) {
				if tt.tasks != nil {
					return tt.tasks, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := ListTask{
				Service: moq,
			}
			sut.ServeHTTP(w, r)
			rsp := w.Result()
			testutil.AsertResponse(t, rsp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
