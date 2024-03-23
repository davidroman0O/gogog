package main

import (
	"testing"

	"github.com/davidroman0O/gogog/db"
	"github.com/davidroman0O/gogog/jobs"
)

func TestJobTest(t *testing.T) {

	if err := jobs.Initialize(); err != nil {
		t.Error(err)
	}

	if err := db.Initialize(
		db.WithDBConfig(db.DBWithMode(db.Memory)),
		db.WithMiddleware[jobs.JobMiddleware](
			db.WithInitializer[jobs.JobMiddleware](jobs.Initializer()),
			db.WithCloser[jobs.JobMiddleware](jobs.Closer()),
			db.WithOnInsert[jobs.JobMiddleware](jobs.OnInsert()),
		),
	); err != nil {
		t.Error(err)
	}

	if err := db.Close(); err != nil {
		t.Error(err)
	}
}
