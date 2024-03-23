package data

import (
	"database/sql"
	"log"

	"github.com/mattn/go-sqlite3"
)

// TODO @droman: implement a task

func init() {
	sql.Register("hook_data",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
					switch op {
					case sqlite3.SQLITE_INSERT:
						log.Println("Notified of insert on db", db, "table", table, "rowid", rowid)
					}
				})
				return nil
			},
		})
}

type Task struct{}
