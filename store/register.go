package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/go-sql-driver/mysql"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	// 時刻情報
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `insert into user (name,password,role,created,modified) values(?,?,?,?,?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}
