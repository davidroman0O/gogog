package jobs

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/mattn/go-sqlite3"
)

/// https://sqldocs.org/sqlite/golang-sqlite/
/// TODO @droman: I will do a simple job system that leverage the sqlite
/// The basics are working, i have to deal with database locks and concurrency. For now it's really dirty and i will have to make some cleaning but it's working!

var global *JobMiddleware
var wasInitialized bool

func Middleware() *JobMiddleware {
	return global
}

func Initialize() error {
	if global == nil {
		global = &JobMiddleware{
			tasks: map[string][]FnTaskCallback{},
		}
	}
	wasInitialized = true
	return nil
}

var jobsTable = `CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    payload JSON,
    status TEXT NOT NULL DEFAULT 'enqueued',
    created_at DATETIME,
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

// A job can contain any payload and can be hard deleted.
// Only the runtime knows how to handle a job type with a callback.
// A job is a stored item that represent a future Task (runtime)
type job[T any] struct {
	ID        int64
	State     State        `json:"status"`
	Type      string       `json:"type"`
	Payload   T            `json:"payload"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

// A work is a pending job to be executed by the runtime
type work interface{}

// task is a pending `TaskCallback` to be executed by the runtime
type task interface{}

type taskCallback[T any] func(data T) error

func On[T any](cb taskCallback[T]) error {
	if !wasInitialized {
		return fmt.Errorf("jobs system was not initialized")
	}
	return global.Add(reflect.TypeFor[T]().Name(), cb)
}

func Do(work work) error {
	if !wasInitialized {
		return fmt.Errorf("jobs system was not initialized")
	}
	return global.Create(work)
}

type FnTaskCallback struct {
	fn reflect.Value
	in reflect.Type
}

type JobMiddleware struct {
	tasks map[string][]FnTaskCallback
	db    *sql.DB
}

func (t *JobMiddleware) Create(w work) error {

	var valueOfWork work = w
	if reflect.TypeOf(w).Kind() == reflect.Ptr {
		valueOfWork = reflect.ValueOf(w).Elem().Interface().(work)
	}

	nameType := reflect.TypeOf(valueOfWork).Name()

	// var prepared *sql.Stmt
	var err error

	tx, err := t.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// fmt.Println("prepare insert into db")

	var dataJson []byte
	if dataJson, err = json.Marshal(valueOfWork); err != nil {
		return err
	}

	// fmt.Println("exec insert into db")
	if _, err = tx.Exec("INSERT INTO jobs (status, type, payload, created_at) VALUES (?, ?, ?, datetime('now'))", Enequeued, nameType, string(dataJson)); err != nil {
		return err
	}

	// fmt.Println("insert db log", t.db)

	return tx.Commit()
}

func (t *JobMiddleware) Add(typeFor string, task task) error {
	taskType := reflect.TypeOf(task)
	if taskType.Kind() != reflect.Func {
		return fmt.Errorf("task must be a function")
	}
	if taskType.NumIn() != 1 || taskType.NumOut() != 1 || taskType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("task function must have one parameter of type 'data' and return 'error'")
	}

	// that's a whole old trick i used in my go-experiments to abstract the functions and their parameters so i can have a lit api
	fn := reflect.ValueOf(task)
	// if fn.Type().NumIn() == 1 {
	// 	// Get the type of the first parameter
	// 	paramType := fn.Type().In(0)

	// 	// Continue with paramType...
	// }

	// if _, ok := t.tasks[typeFor]; !ok {
	// 	t.tasks[typeFor] = []task{}
	// }
	t.tasks[typeFor] = append(t.tasks[typeFor], FnTaskCallback{fn: fn, in: fn.Type().In(0)})
	// t.tasks = append(t.tasks, task)
	return nil
}

func (t *JobMiddleware) prepareTables(db *sql.DB) error {

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	// defer tx.Rollback()

	data, err := tx.Exec(jobsTable)
	fmt.Println(data.RowsAffected())
	fmt.Println(data.LastInsertId())
	fmt.Println(err)

	{
		// Execute the query
		query := `SELECT name FROM sqlite_schema WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%' ORDER BY 1;`
		rows, err := tx.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Print the results
		fmt.Println("Tables and Views:")
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Fatal(err)
			}
			fmt.Println(name)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return tx.Commit()
}

func (t *JobMiddleware) OnInit(db *sql.DB) error {
	log.Println("Jobs middleware initialized")
	t.db = db
	return t.prepareTables(db)
}

func (t *JobMiddleware) OnClose(db *sql.DB) error {
	log.Println("Jobs middleware closed")
	return nil
}

func (t *JobMiddleware) OnInsert(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error {
	// log.Printf("Jobs Insert operation on table %s with rowid %d", table, rowid)
	var err error

	// fmt.Println("lock before beginx tx hook")

	rows, err := conn.Query(`SELECT id, type, payload, status, strftime('%Y-%m-%d %H:%M:%S', created_at) FROM jobs WHERE id = ?`, []driver.Value{rowid})
	if err != nil {
		return err
	}
	defer rows.Close()

	var results map[string][]job[any] = make(map[string][]job[any])

	dest := make([]driver.Value, 5)
	for {
		err = rows.Next(dest)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		id, _ := dest[0].(int64)
		typeJob, _ := dest[1].(string)
		payload, _ := dest[2].(string)
		status, _ := dest[3].(string)
		created_at, _ := dest[4].(string)

		// fmt.Println(id, typeJob, payload, status, dest[4])
		createdAt, err := time.Parse("2006-01-02 15:04:05", created_at)
		if err != nil {
			return err // Handle parsing error
		}

		// find those interested in this type of job
		for _, receivers := range t.tasks[typeJob] {

			row := job[any]{
				ID:        id,
				State:     State(status),
				Type:      typeJob,
				CreatedAt: createdAt,
			}

			paramInstancePtr := reflect.New(receivers.in).Interface()

			err = json.Unmarshal([]byte(payload), paramInstancePtr)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// pp.Println(paramInstancePtr)

			row.Payload = paramInstancePtr

			results[typeJob] = append(results[typeJob], row)
		}

	}

	// pp.Println(results)

	for key, rows := range results {
		if t.tasks[key] == nil {
			continue
		}
		for _, receiver := range t.tasks[key] {

			for _, v := range rows {

				// Convert paramInstancePtr (which is an interface{}) back to reflect.Value
				// Note: If your function expects a value instead of a pointer, you might need to adjust this
				// pp.Println(v.Payload)
				paramValue := reflect.ValueOf(v.Payload)

				results := receiver.fn.Call([]reflect.Value{
					paramValue.Elem(),
				})
				if !results[0].IsNil() {
					return results[0].Interface().(error)
				}
			}

		}
	}

	return nil
}

func (t *JobMiddleware) OnUpdate(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error {
	log.Printf("Jobs Update operation on table %s with rowid %d", table, rowid)
	return nil
}

func (t *JobMiddleware) OnDelete(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error {
	log.Printf("Jobs Delete operation on table %s with rowid %d", table, rowid)
	return nil
}
