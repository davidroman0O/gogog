package jobs

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/mattn/go-sqlite3"
)

type JobData struct {
	Msg string
}

func TestSimple(t *testing.T) {
	var db *sql.DB
	var err error

	sql.Register(
		"sqlite3-extended",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// make so something with `conn`?
				conn.RegisterUpdateHook(
					func(op int, db string, table string, rowid int64) {
						switch op {

						case sqlite3.SQLITE_INSERT:
							global.OnInsert(conn, db, table, rowid)

						case sqlite3.SQLITE_UPDATE:
							global.OnUpdate(conn, db, table, rowid)

						case sqlite3.SQLITE_DELETE:
							global.OnDelete(conn, db, table, rowid)

						}
					},
				)
				return nil
			}})

	// TODO: I have a bug, when it's full memory, it's not working `?mode=memory`
	// if db, err = sql.Open("sqlite3-extended", "file:something.db?_journal_mode=WAL&_mutex=full"); err != nil {
	// 	t.Error(err)
	// }
	var connectionString = ":memory:?cache=shared"

	// if db, err = sql.Open("sqlite3-extended", "file::memory:?cache=shared&_busy_timeout=5000"); err != nil {
	// if db, err = sql.Open("sqlite3-extended", "file:memdb1?mode=memory&cache=shared&_busy_timeout=5000"); err != nil {
	if db, err = sql.Open("sqlite3-extended", connectionString); err != nil {
		t.Error(err)
	}
	// db.Exec("PRAGMA journal_mode=WAL;")

	// Sqlite cannot handle concurrent writes,
	// so we limit sqlite to one connection.
	// https://github.com/mattn/go-sqlite3/issues/274
	// db.SetMaxOpenConns(1)

	// if _, err := db.Exec("ATTACH DATABASE 'file::memory:?cache=shared' AS inmem"); err != nil {
	// 	panic(err)
	// }
	// db.SetMaxOpenConns(1)

	// if db, err = sql.Open("sqlite3-extended", "file::memory:?_journal_mode=MEMORY&mode=memory&cache=shared&_mutex=full"); err != nil {
	// 	t.Error(err)
	// }

	if err := Initialize(); err != nil {
		t.Error(err)
	}

	if err := global.OnInit(db); err != nil {
		t.Error(err)
	}

	if err := On(func(data JobData) error {
		fmt.Println("triggered", data)
		return nil
	}); err != nil {
		t.Error(err)
	}

	if err := Do(JobData{Msg: "hello"}); err != nil {
		t.Error(err)
	}

	// Execute the query
	query := `SELECT name FROM sqlite_schema WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%' ORDER BY 1;`
	rows, err := db.Query(query)
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
