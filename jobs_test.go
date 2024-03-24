package main

import (
	"fmt"
	"testing"

	"github.com/davidroman0O/gogog/db"
	"github.com/davidroman0O/gogog/jobs"
)

type JobData struct {
	Msg string
}

func TestJobTest(t *testing.T) {

	if err := jobs.Initialize(); err != nil {
		t.Error(err)
	}

	middleware := jobs.Middleware()

	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := db.Initialize(
		db.WithDBConfig(db.DBWithMode(db.Memory), db.DBWithCacheShared()),
		db.WithMiddleware(middleware),
	); err != nil {
		t.Error(err)
	}

	if err := jobs.On(func(data JobData) error {
		fmt.Println("triggered", data)
		return nil
	}); err != nil {
		t.Error(err)
	}

	if err := jobs.Do(JobData{Msg: "hello"}); err != nil {
		t.Error(err)
	}
}
