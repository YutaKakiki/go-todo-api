package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YutaKakiki/go-todo-api/testutil"
	"github.com/go-playground/validator/v10"
)

func TestLogin_ServerHTTP(t *testing.T) {
	// モックが返す型
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		moq     moq
		want    want
	}{
		"ok": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/login/ok_rsp.json.golden",
			},
		},
		"badReq": {
			reqFile: "testdata/login/bad_req.json.golden",
			// Loginメソッドの前にバリデーションに引っかかるのでmockの値を設定していない
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/login/bad_rsp.json.golden",
			},
		},
		"internalServerErr": {
			reqFile: "testdata/login/ok_req.json.golden",
			// Loginメソッドの挙動
			moq: moq{
				err: errors.New("error from mock"),
			},
			want: want{
				status:  http.StatusInternalServerError,
				rspFile: "testdata/login/internal_server_error.json.golden",
			},
		},
	}
	for nn, tt := range tests {

		tt := tt
		t.Run(nn, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet, "/login", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)
			// サービス層のLoginメソッドをモック化
			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, name, pw string) (string, error) {
				return tt.moq.token, tt.moq.err
			}

			sut := Login{
				Service:   moq,
				Validator: validator.New(),
			}
			// ServeHTTP
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AsertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})

	}

}
