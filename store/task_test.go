package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YutaKakiki/go-todo-api/clock"
	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	if _, err := con.ExecContext(ctx, "delete from task;"); err != nil {
		t.Logf("failed to init task:%v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title:    "want task 1",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	// ListTaskのテストで期待するデータを流し込む
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
			VALUES
			    (?, ?, ?, ?),
			    (?, ?, ?, ?),
			    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}

func TestRepository_ListTask(t *testing.T) {
	ctx := context.Background()
	// テスト用のDB
	// トランザクション開始：このテストケースのみのテーブルの状態に持っていく
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	// テストが終わったらロールバック
	t.Cleanup(func() { tx.Rollback() })

	wants := prepareTasks(ctx, t, tx)
	sut := &Repository{}
	gots, err := sut.ListTask(ctx, tx)
	if err != nil {
		t.Fatalf("unexped error:%v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs:(-got +want)\n%s", d)
	}
}

// モックDBでテスト
// RDBMSに依存しないで書ける
func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}
	// 	db: モックデータベースを表す *sql.DB 型のオブジェクト。
	// mock: SQL モックの操作を定義するための sqlmock.Sqlmock 型のオブジェクト。
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	// モックDBは、以下の操作を期待する
	mock.ExpectExec(
		`insert into task \(title, status, created, modified\) values \(\?, \?, \?, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
		// モックデータベースがこのクエリに対して返す結果
		// 挿入された行の ID (wantID) と、影響を受けた行数（ここでは 1）
		WillReturnResult(sqlmock.NewResult(wantID, 1))
	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error,but got:%v", err)
	}
	// 期待したクエリが発行されたのか
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
