package db

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	logger "github.com/davidroman0O/gogog/db/middlewares"
)

/// TODO @droman: test all options
/// TODO @droman: i need to find to  completely remove the sqlite3 context between tests to avoid the "double registry" issue

func TestOpenCloseMemory(t *testing.T) {

	if err := Initialize(
		WithDBConfig(DBWithMode(Memory))); err != nil {
		t.Error(err)
	}

	didPing := false

	if err := Do(func(db *sql.DB) error {
		if err := db.Ping(); err != nil {
			return err
		}
		slog.Info("pinged")
		didPing = true
		return nil
	}); err != nil {
		t.Error(err)
	}

	if !didPing {
		t.Error("did not ping")
	}

	if err := Close(); err != nil {
		t.Error(err)
	}
}

type MyData struct {
	ID        int64
	Natural   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	Data      map[string]interface{}
}

func TestCreateTableMemory(t *testing.T) {

	logger := &logger.LoggerMiddleware{}

	if err := Initialize(
		WithMiddleware(logger),
		WithDBConfig(
			DBWithMode(Memory))); err != nil {
		t.Error(err)
	}

	// Let's test a basic table with a few columns
	if err := Do(func(db *sql.DB) error {
		_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS mytable (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	natural TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NULL,
	deleted_at DATETIME NULL,
	data JSON
);
		`)
		return err
	}); err != nil {
		t.Error(err)
	}

	if err := Do(func(db *sql.DB) error {
		_, err := db.Exec("INSERT INTO mytable (natural, data) VALUES (?, ?)", "Example 1", `{"key": "value", "numbers": [1, 2, 3]}`)
		return err
	}); err != nil {
		t.Error(err)
	}

	myData := MyData{
		Natural: "Example 2",
		Data: map[string]interface{}{
			"key":     "value",
			"numbers": []int{1, 2, 3},
		},
		CreatedAt: time.Now(),
	}

	dataJSON, err := json.Marshal(myData.Data)
	if err != nil {
		t.Error(err)
	}

	if err := Do(func(db *sql.DB) error {
		_, err := db.Exec(`
        INSERT INTO mytable (natural, data, created_at)
        VALUES (?, ?, ?)
    `, myData.Natural, dataJSON, myData.CreatedAt)
		return err
	}); err != nil {
		t.Error(err)
	}

	var results []MyData
	if err := Do(func(db *sql.DB) error {
		rows, err := db.Query("SELECT id, natural, created_at, updated_at, deleted_at, data FROM mytable")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var row MyData
			var dataJSON []byte

			err = rows.Scan(&row.ID, &row.Natural, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt, &dataJSON)
			if err != nil {
				return err
			}

			err = json.Unmarshal(dataJSON, &row.Data)
			if err != nil {
				return err
			}

			results = append(results, row)
		}

		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Error(err)
	}

	// pp.Println(results)

	if len(results) == 0 {
		t.Error("no results")
	}

	slog.Info("closing")

	if err := Close(); err != nil {
		t.Error(err)
	}
}
