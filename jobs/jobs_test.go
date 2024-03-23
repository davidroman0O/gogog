package jobs

import (
	"database/sql"
	"fmt"
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
				conn.RegisterUpdateHook(
					func(op int, db string, table string, rowid int64) {
						switch op {

						case sqlite3.SQLITE_INSERT:
							global.OnInsert(db, table, rowid)
							// slog.Info("hook OnInsert", slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))

						case sqlite3.SQLITE_UPDATE:
							global.OnUpdate(db, table, rowid)
							// slog.Info("hook OnUpdate", slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))

						case sqlite3.SQLITE_DELETE:
							global.OnDelete(db, table, rowid)
							// slog.Info("hook OnDelete", slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))
						}
					},
				)

				return nil
			}})

	if db, err = sql.Open("sqlite3-extended", "file:something.db?mode=memory"); err != nil {
		t.Error(err)
	}

	if err := Initialize(); err != nil {
		t.Error(err)
	}

	// simulate db initialization
	global.Initialize(func(cb func(db *sql.DB) error) error {
		return cb(db)
	})

	if err := On(func(data JobData) error {
		fmt.Println("triggered", data)
		return nil
	}); err != nil {
		t.Error(err)
	}

	if err := Do(JobData{Msg: "hello"}); err != nil {
		t.Error(err)
	}
}
