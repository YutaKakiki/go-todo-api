package store

import (
	"context"

	"github.com/YutaKakiki/go-todo-api/entity"
)

// dbは、Queryer型である必要がある
// (*sqlx.DBとか)
func (r *Repository) ListTask(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `select id,title,status,created,modified from task;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `insert into task (title, status, created, modified) values (?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.Created, t.Modified)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil

}
