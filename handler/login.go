package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// リクエストボディの型
	var body struct {
		UserName string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	err = l.Validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// サービス層に処理を引き継ぐ
	jwt, err := l.Service.Login(ctx, body.UserName, body.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt,
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
