package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/YutaKakiki/go-todo-api/clock"
	"github.com/YutaKakiki/go-todo-api/entity"
	"github.com/YutaKakiki/go-todo-api/testutil"
	"github.com/YutaKakiki/go-todo-api/testutil/fixture"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()
	u := fixture.User(nil)
	result, err := db.ExecContext(ctx, `insert into user (name,password,role,created,modified) values (?,?,?,?,?)`, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		t.Fatalf("insert user :%v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get user_id,:%v", err)
	}
	return entity.UserID(id)

}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()
	if _, err := con.ExecContext(ctx, "delete from task;"); err != nil {
		t.Logf("failed to init task:%v", err)
	}
	userID := prepareUser(ctx, t, con)
	otherID := prepareUser(ctx, t, con)
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID:   userID,
			Title:    "want task 1",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		{
			UserID: userID,
			Title:  "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserID:   otherID,
			Title:    "not want task",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		wants[1],
	}
	// ListTaskのテストで期待するデータを流し込む
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (user_id,title, status, created, modified)
			VALUES
			    (?, ?, ?, ?, ?),
			    (?, ?, ?, ?, ?),
			    (?, ?, ?, ?, ?);`,
		tasks[0].UserID, tasks[0].Title, tasks[0].Status, tasks[0].Created, tasks[0].Modified,
		tasks[1].UserID, tasks[1].Title, tasks[1].Status, tasks[1].Created, tasks[1].Modified,
		tasks[2].UserID, tasks[2].Title, tasks[2].Status, tasks[2].Created, tasks[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
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

	wantID, wants := prepareTasks(ctx, t, tx)
	sut := &Repository{}
	gots, err := sut.ListTask(ctx, tx, wantID)
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
		UserID:   33,
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
		`insert into task \(user_id, title, status, created, modified\) values \(\?, \?, \?, \?, \?\)`,
	).WithArgs(okTask.UserID, okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
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
