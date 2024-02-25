package data

import (
	"fmt"

	"github.com/asdine/storm/v3"
)

var accountsDB = "accounts.db"

func initializeAccounts(prefix string) error {
	if prefix != "" {
		accountsDB = fmt.Sprintf("%v/%v", prefix, accountsDB)
	}
	var err error
	if dbAccounts, err = storm.Open(accountsDB); err != nil {
		return err
	}
	return nil
}

var dbAccounts *storm.DB

func Initialize(prefix string) error {

	if err := initializeAccounts(prefix); err != nil {
		return err
	}

	return nil
}

func Close() {
	if dbAccounts != nil {
		dbAccounts.Close()
	}
}

func AccountsDB() *storm.DB {
	return dbAccounts
}
