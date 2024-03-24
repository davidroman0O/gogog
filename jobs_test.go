package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

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

	triggered := 0
	now := time.Now()

	if err := jobs.On(func(data JobData) error {
		// fmt.Println("triggered", data.Msg)
		triggered++
		return nil
	}); err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup

	// if max connection is one, we can do multi goroutines
	for range 2 {
		wg.Add(1)
		go func() {
			for range 100000 {
				if err := jobs.Do(JobData{Msg: "hello"}); err != nil {
					t.Error(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("elapsed", time.Since(now), triggered)

}
