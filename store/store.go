package store

import "errors"

// import (
// 	"errors"

// 	"github.com/YutaKakiki/go-todo-api/entity"
// )

var (
	// Tasks       = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}
	ErrNotFound = errors.New("not found")
)

// type TaskStore struct {
// 	LastID entity.TaskID
// 	Tasks  map[entity.TaskID]*entity.Task
// }

// func (ts *TaskStore) Add(t *entity.Task) (entity.TaskID, error) {
// 	// 最後のレコードのIDをインクリメント
// 	ts.LastID++
// 	// ＋１したIDを設定
// 	t.ID = ts.LastID
// 	// TaskStore構造体に代入
// 	ts.Tasks[t.ID] = t
// 	return t.ID, nil
// }

// func (ts *TaskStore) Get(id entity.TaskID) (*entity.Task, error) {
// 	if ts, ok := ts.Tasks[id]; ok {
// 		return ts, nil
// 	}
// 	return nil, ErrNotFound
// }

// func (ts *TaskStore) GetAll() entity.Tasks {
// 	tasks := make([]*entity.Task, len(ts.Tasks))
// 	for i, t := range ts.Tasks {
// 		tasks[i-1] = t
// 	}
// 	return tasks
// }
