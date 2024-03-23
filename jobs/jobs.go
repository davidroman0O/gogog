package jobs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"
)

/// https://sqldocs.org/sqlite/golang-sqlite/
/// TODO @droman: I will do a simple job system that leverage the sqlite

var jobsTable = `
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    payload JSON,
    status TEXT NOT NULL DEFAULT 'enqueued',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL,
    error TEXT NULL
);
`

type State string

var (
	Enequeued State = "enqueued"
	Pending   State = "pending"
	Running   State = "running"
	Failed    State = "failed"
)

// A Job can contain any payload and can be hard deleted.
// Only the runtime knows how to handle a job type with a callback.
// A Job is a stored item that represent a future Task (runtime)
type Job[T any] struct {
	ID        int64
	State     State
	Type      string
	Payload   T
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// Task is a pending `TaskCallback` to be executed by the runtime
type Task interface{}

type TaskCallback[T any] func(data T) error

func On[T any](cb TaskCallback[T]) error {
	if !wasInitialized {
		return fmt.Errorf("jobs system was not initialized")
	}
	return global.Add(cb)
}

func Do[T any](data T) error {

	newJob := Job[T]{
		State:     Enequeued,
		Type:      reflect.TypeFor[T]().Name(),
		Payload:   data,
		CreatedAt: time.Now(),
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return global.do(func(db *sql.DB) error {
		_, err := db.Exec(`
        INSERT INTO jobs (status, type, payload, created_at)
        VALUES (?, ?, ?, ?)
    	`, newJob.State, newJob.Type, dataJSON, newJob.CreatedAt)
		return err
	})
}

type JobMiddleware struct {
	tasks []Task
	mu    sync.Mutex
	do    func(cb func(db *sql.DB) error) error
}

func (t *JobMiddleware) Add(task Task) error {
	taskType := reflect.TypeOf(task)
	if taskType.Kind() != reflect.Func {
		return fmt.Errorf("Task must be a function")
	}
	if taskType.NumIn() != 1 || taskType.NumOut() != 1 || taskType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("Task function must have one parameter of type 'data' and return 'error'")
	}
	t.tasks = append(t.tasks, task)
	return nil
}

func (m *JobMiddleware) Initialize(pass func(cb func(db *sql.DB) error) error) error {
	if pass == nil {
		return fmt.Errorf("pass function is required")
	}
	m.do = pass
	m.mu.Lock()

	if err := pass(func(db *sql.DB) error {
		_, err := db.Exec(jobsTable)
		return err
	}); err != nil {
		return err
	}

	m.mu.Unlock()
	slog.Info("job system initialized")
	return nil
}

func (m *JobMiddleware) Close() error {
	return nil
}

func (m *JobMiddleware) OnInsert(db string, table string, rowid int64) error {
	slog.Info("job middleware insert hook", slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))
	return nil
}

func (m *JobMiddleware) OnUpdate(db string, table string, rowid int64) error {
	slog.Info("job middleware update hook", slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))
	return nil
}

func (m *JobMiddleware) OnDelete(db string, table string, rowid int64) error {
	slog.Info("job middleware delete hook", slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))
	return nil
}

var global *JobMiddleware
var wasInitialized bool

func Initializer() func(pass func(cb func(db *sql.DB) error) error) error {
	return global.Initialize
}

func Closer() func() error {
	return global.Close
}

func OnInsert() func(db string, table string, rowid int64) error {
	return global.OnInsert
}

func OnUpdate() func(db string, table string, rowid int64) error {
	return global.OnUpdate
}

func OnDelete() func(db string, table string, rowid int64) error {
	return global.OnDelete
}

func Initialize() error {
	if global == nil {
		global = &JobMiddleware{
			tasks: []Task{},
		}
	}
	wasInitialized = true
	return nil
}
