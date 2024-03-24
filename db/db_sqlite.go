package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/mattn/go-sqlite3"
)

var db *sql.DB
var config *initConfig

type DoFn func(db *sql.DB) error
type PassDoFn func(cb DoFn) error

func Do(cb DoFn) error {
	return cb(db)
}

type initConfig struct {
	db                *dbConfig
	middlewareManager *MiddlewareManager
}

type initOpts func(*initConfig) error

func WithDBConfig(opts ...dbOption) initOpts {
	return func(config *initConfig) error {
		config.db = NewSettingConfig(opts...)
		return nil
	}
}

func WithDBMemory() initOpts {
	return func(config *initConfig) error {
		config.db = NewSettingConfig(
			DBWithMode(Memory),
		)
		return nil
	}
}

func WithMiddleware(middleware Middleware) initOpts {
	return func(ic *initConfig) error {
		if reflect.TypeOf(middleware).Kind() != reflect.Ptr {
			return fmt.Errorf("middleware must be a pointer to a struct")
		}
		ic.middlewareManager.Register(middleware)
		return nil
	}
}

func Initialize(opts ...initOpts) error {
	if config == nil {
		config = &initConfig{
			middlewareManager: &MiddlewareManager{},
		}
		for _, opt := range opts {
			if err := opt(config); err != nil {
				return err
			}
		}
	}

	var err error

	var connectionString string
	if connectionString, err = ConnectionString(config.db); err != nil {
		return err
	}

	slog.Info("connection string ", slog.String("value", connectionString))

	sql.Register(
		config.db.name,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// register callback
				conn.RegisterUpdateHook(
					func(op int, db string, table string, rowid int64) {

						switch op {

						case sqlite3.SQLITE_INSERT:
							if err := config.middlewareManager.RunOnInsert(conn, db, table, rowid); err != nil {
								// TODO: Handle error - i have no idea how
								slog.Error("SQLITE_INSERT error: %v", err)
							}

						case sqlite3.SQLITE_UPDATE:
							if err := config.middlewareManager.RunOnUpdate(conn, db, table, rowid); err != nil {
								// TODO: Handle error - i have no idea how
								slog.Error("SQLITE_UPDATE error: %v", err)
							}

						case sqlite3.SQLITE_DELETE:
							if err := config.middlewareManager.RunOnDelete(conn, db, table, rowid); err != nil {
								// TODO: Handle error - i have no idea how
								slog.Error("SQLITE_DELETE error: %v", err)
							}

						}
					},
				)

				return nil
			}})

	if db, err = sql.Open(config.db.name, connectionString); err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Initialize middlewares
	if err := config.middlewareManager.RunOnInit(db); err != nil {
		return err
	}

	return nil
}

func Close() error {
	if db != nil {

		// Close middlewares
		if err := config.middlewareManager.RunOnClose(db); err != nil {
			return err
		}

		// Check if the file exists
		if _, err := os.Stat(config.db.filePath); err == nil {
			// Remove the file
			if err := os.Remove(config.db.filePath); err != nil {
				return err
			}
		}
		db.Close()
	}

	return nil
}

type Middleware interface {
	OnInit(db *sql.DB) error
	OnClose(db *sql.DB) error
	OnInsert(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error
	OnUpdate(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error
	OnDelete(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error
}

type MiddlewareManager struct {
	middlewares []Middleware
}

func (mm *MiddlewareManager) Register(middleware Middleware) {
	mm.middlewares = append(mm.middlewares, middleware)
}

func (mm *MiddlewareManager) RunOnInit(db *sql.DB) error {
	for _, middleware := range mm.middlewares {
		if err := middleware.OnInit(db); err != nil {
			return err
		}
	}
	return nil
}

func (mm *MiddlewareManager) RunOnClose(db *sql.DB) error {
	for _, middleware := range mm.middlewares {
		if err := middleware.OnClose(db); err != nil {
			return err
		}
	}
	return nil
}

func (mm *MiddlewareManager) RunOnInsert(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error {
	for _, middleware := range mm.middlewares {
		if err := middleware.OnInsert(conn, db, table, rowid); err != nil {
			return err
		}
	}
	return nil
}

func (mm *MiddlewareManager) RunOnUpdate(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error {
	for _, middleware := range mm.middlewares {
		if err := middleware.OnUpdate(conn, db, table, rowid); err != nil {
			return err
		}
	}
	return nil
}

func (mm *MiddlewareManager) RunOnDelete(conn *sqlite3.SQLiteConn, db string, table string, rowid int64) error {
	for _, middleware := range mm.middlewares {
		if err := middleware.OnDelete(conn, db, table, rowid); err != nil {
			return err
		}
	}
	return nil
}
