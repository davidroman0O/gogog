package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var config *initConfig
var connectionString string

func Do(cb func(db *sql.DB) error) error {
	return cb(db)
}

type initConfig struct {
	db *dbConfig
}

type initOpts func(*initConfig)

func WithDBConfig(opts ...dbOption) initOpts {
	return func(config *initConfig) {
		config.db = NewSettingConfig(opts...)
	}
}

func Initialize(opts ...initOpts) error {
	if config == nil {
		config = &initConfig{}
		for _, opt := range opts {
			opt(config)
		}
	}

	var err error
	var connectionString string

	if connectionString, err = ConnectionString(config.db); err != nil {
		return err
	}

	if db, err = sql.Open("sqlite3", connectionString); err != nil {
		return err
	}

	return nil
}

func Close() error {
	if db != nil {
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
