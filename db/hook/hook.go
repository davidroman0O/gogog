package hook

import (
	"database/sql"
	"log"
	"log/slog"

	"github.com/mattn/go-sqlite3"
)

func Hook(name string) {

	sql.Register(
		name,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {

				slog.Info("hooked")

				// register callback
				conn.RegisterUpdateHook(
					func(op int, db string, table string, rowid int64) {
						slog.Info("hook", slog.Any("op", op), slog.Any("db", db), slog.Any("table", table), slog.Any("rowid", rowid))
						switch op {
						case sqlite3.SQLITE_INSERT:
							log.Println("Notified of insert on db", db, "table", table, "rowid", rowid)
						case sqlite3.SQLITE_UPDATE:
							log.Println("Notified of update on db", db, "table", table, "rowid", rowid)
						case sqlite3.SQLITE_DELETE:
							log.Println("Notified of delete on db", db, "table", table, "rowid", rowid)
						}
					},
				)

				return nil
			}})
}
