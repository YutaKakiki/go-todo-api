package service

import (
	"context"
	"fmt"

	"github.com/YutaKakiki/go-todo-api/store"
)

type Login struct {
	DB store.Queryer
	// ユーザーを得るstore(repository)とのインターフェース
	Repo UserGetter
	// auth.JWTerとのインターフェース
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {

	// nameを使ってDBからUserを取得
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	// DB内にあるハッシュ化されたパスワードと入力したパスワードを比較
	err = u.CompairePassword(pw)
	if err != nil {
		return "", fmt.Errorf("wrong password : %w", err)
	}
	// jwtを発行
	jwt, err := l.TokenGenerator.GenerateToken(ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to gen JWT : %w", err)
	}
	return string(jwt), nil
}
