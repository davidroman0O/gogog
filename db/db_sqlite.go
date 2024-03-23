package db

import (
	"database/sql"
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
	db          *dbConfig
	middlewares map[reflect.Type]Middleware
}

type initOpts func(*initConfig) error

func WithDBConfig(opts ...dbOption) initOpts {
	return func(config *initConfig) error {
		config.db = NewSettingConfig(opts...)
		return nil
	}
}

type Middleware struct {
	Initializer func(pass func(cb func(db *sql.DB) error) error) error
	Closer      func() error
	OnInsert    func(db string, table string, rowid int64) error
	OnUpdate    func(db string, table string, rowid int64) error
	OnDelete    func(db string, table string, rowid int64) error
}

type middlewareOptions func(*Middleware) error

func WithInitializer[T any](Initializer func(pass func(cb func(db *sql.DB) error) error) error) middlewareOptions {
	return func(middleware *Middleware) error {
		middleware.Initializer = Initializer
		return nil
	}
}

func WithCloser[T any](closer func() error) middlewareOptions {
	return func(middleware *Middleware) error {
		middleware.Closer = closer
		return nil
	}
}
func WithOnInsert[T any](hook func(db string, table string, rowid int64) error) middlewareOptions {
	return func(middleware *Middleware) error {
		middleware.OnInsert = hook
		return nil
	}
}

func WithOnUpdate[T any](hook func(db string, table string, rowid int64) error) middlewareOptions {
	return func(middleware *Middleware) error {
		middleware.OnUpdate = hook
		return nil
	}
}

func WithOnDelete[T any](hook func(db string, table string, rowid int64) error) middlewareOptions {
	return func(middleware *Middleware) error {
		middleware.OnDelete = hook
		return nil
	}
}

func WithMiddleware[T any](opts ...middlewareOptions) initOpts {
	return func(config *initConfig) error {
		middleware := &Middleware{}
		for _, v := range opts {
			if err := v(middleware); err != nil {
				return err
			}
		}
		config.middlewares[reflect.TypeFor[T]()] = *middleware
		return nil
	}
}

// func Find[T Middleware]() (*T, bool) {
// 	for _, middleware := range config.middlewares {
// 		if m, ok := (middleware).(T); ok {
// 			return &m, true
// 		}
// 	}
// 	return nil, false
// }

func Initialize(opts ...initOpts) error {
	if config == nil {
		config = &initConfig{
			middlewares: map[reflect.Type]Middleware{},
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

	sql.Register(
		config.db.name,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// register callback
				conn.RegisterUpdateHook(
					func(op int, db string, table string, rowid int64) {

						switch op {

						case sqlite3.SQLITE_INSERT:
							for _, middleware := range config.middlewares {
								if middleware.OnInsert != nil {
									middleware.OnInsert(db, table, rowid) // TODO @droman: how the fuck i will managed that?
								}
							}

						case sqlite3.SQLITE_UPDATE:
							for _, middleware := range config.middlewares {
								if middleware.OnUpdate != nil {
									middleware.OnUpdate(db, table, rowid) // TODO @droman: how the fuck i will managed that?
								}
							}

						case sqlite3.SQLITE_DELETE:
							for _, middleware := range config.middlewares {
								if middleware.OnDelete != nil {
									middleware.OnDelete(db, table, rowid) // TODO @droman: how the fuck i will managed that?
								}
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

	for _, middleware := range config.middlewares {
		if err := middleware.Initializer(func(cb func(db *sql.DB) error) error { return cb(db) }); err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		for _, middleware := range config.middlewares {
			if err := middleware.Closer(); err != nil {
				return err
			}
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
