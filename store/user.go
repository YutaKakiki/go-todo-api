package store

import (
	"context"

	"github.com/YutaKakiki/go-todo-api/entity"
)

func (r *Repository) GetUser(ctx context.Context, db Queryer, name string) (*entity.User, error) {
	u := &entity.User{}
	sql := `select id,name,password,role,created,modified from user where name = ?`
	if err := db.GetContext(ctx, u, sql, name); err != nil {
		return nil, err
	}
	return u, nil
}
