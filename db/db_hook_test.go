package db

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"log/slog"
	"testing"

	"github.com/mattn/go-sqlite3"
)

func TestHooks(t *testing.T) {
	var db *sql.DB
	var err error

	var connectionString = ":memory:?cache=shared"
	var deferConn *sqlite3.SQLiteConn
	sql.Register(
		"sqlite3-extended",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// Defer the registration of the update hook until after the table has been created
				deferConn = conn
				return nil
			}})

	if db, err = sql.Open("sqlite3-extended", connectionString); err != nil {
		t.Error(err)
	}

	db.SetMaxOpenConns(1)

	slog.Info("Creating table")
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS jobs (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                status TEXT NOT NULL DEFAULT 'enqueued'
            );
        `); err != nil {
		t.Error(err)
	}
	slog.Info("Created table")

	defer deferConn.Close()

	deferConn.RegisterUpdateHook(
		func(op int, dbName string, table string, rowid int64) {
			switch op {
			case sqlite3.SQLITE_INSERT, sqlite3.SQLITE_UPDATE, sqlite3.SQLITE_DELETE:
				slog.Info("hooked")
				slog.Info("selecting jobs")

				rows, err := deferConn.Query("SELECT * FROM jobs", nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer rows.Close()

				dest := make([]driver.Value, 2)
				for {
					err = rows.Next(dest)
					if err == io.EOF {
						break
					} else if err != nil {
						t.Error(err)
						return
					}
					id, _ := dest[0].(int64)
					status, _ := dest[1].(string)
					slog.Info("row", slog.Any("id", id), slog.Any("status", status))
				}

				slog.Info("selected jobs")
			}
		},
	)

	slog.Info("Insertion job")
	if _, err := db.Exec("INSERT INTO jobs (status) VALUES (?);", "enqueued"); err != nil {
		t.Error(err)
	}
	slog.Info("Inserted job")
}

/// hooks works with file `file:thing?_journal_mode=WAL&mode=rwc`
