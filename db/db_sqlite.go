package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var connectionString string

type FileMode string

const (
	ReadWrite FileMode = "rw"
	Memory    FileMode = "memory"
)

type dbConfig struct {
	file FileMode
}

type dbOption func(*dbConfig)

func WithFileMode(fileMode FileMode) dbOption {
	return func(config *dbConfig) {
		config.file = fileMode
	}
}

func NewSettingConfig(options ...dbOption) *dbConfig {
	config := &dbConfig{}
	for _, option := range options {
		option(config)
	}
	return config
}

func getEnvDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
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
	config := initConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	params := map[string]string{
		"file":  fmt.Sprintf("%v:%v", getEnvDefault("DB_PATH", "./tmp/"), getEnvDefault("DB_NAME", "gogog.db")),
		"mode":  getEnvDefault("DB_MODE", "rw"),
		"_sync": getEnvDefault("DB_SYNC", "full"), // we like it safe here
	}

	connectionString = fmt.Sprintf("%v?mode=%v", params["file"], params["mode"])

	var err error
	if db, err = sql.Open("sqlite3", connectionString); err != nil {
		return err
	}

	return nil
}

func Clear() error {
	if db != nil {
		os.Remove("./foo.db")
		db.Close()
	}
	return nil
}
