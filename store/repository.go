package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/YutaKakiki/go-todo-api/clock"
	"github.com/YutaKakiki/go-todo-api/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	// ErrCodeMySQLDuplicateEntry はMySQL系のDUPLICATEエラーコード
	// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	// Error number: 1062; Symbol: ER_DUP_ENTRY; SQLSTATE: 23000
	ErrCodeMySQLDuplicateEntry = 1062
)

var (
	ErrAlreadyEntry = errors.New("duplicate entry")
)

func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		))
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	// データベース接続が正常に動作しているか確認
	// エラーならDBをクローズ
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { db.Close() }, err
	}
	// sqlx.DB にラップする
	xdb := sqlx.NewDb(db, "mysql")
	// コネクションが確立されたらクローズ
	return xdb, func() { db.Close() }, nil
}

// 具体的な型 (sqlx.DB や sqlx.Tx) に直接依存するのではなく、インターフェース（抽象的な定義）に依存させる

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	// １つのクエリ
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	// 複数のクエリ
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

// *sqlx.DBがそれぞれのインターフェースを満たしているのか
// *sqlx.DBもまた上のインターフェースに依存しているため
var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

type Repository struct {
	Clocker clock.Clocker
}
