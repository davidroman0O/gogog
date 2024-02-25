package data

import (
	"errors"
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/davidroman0O/gogog/types"
)

var accountsDB = "accounts.db"
var gamesDB = "games.db"

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

func initializeGames(prefix string) error {
	if prefix != "" {
		gamesDB = fmt.Sprintf("%v/%v", prefix, gamesDB)
	}
	var err error
	if dbGames, err = storm.Open(gamesDB); err != nil {
		return err
	}
	return nil
}

var dbAccounts *storm.DB
var dbGames *storm.DB

func Initialize(prefix string) error {

	if err := initializeAccounts(prefix); err != nil {
		return err
	}

	if err := initializeGames(prefix); err != nil {
		return err
	}

	return nil
}

func Close() {
	if dbAccounts != nil {
		dbAccounts.Close()
	}
	if dbGames != nil {
		dbGames.Close()
	}
}

func AccountsDB() *storm.DB {
	return dbAccounts
}

func GamesDB() *storm.DB {
	return dbGames
}

func CountAccounts() (int, error) {
	return dbAccounts.Count(&types.Account{})
}

func CountGames() (int, error) {
	return dbGames.Count(&types.Account{})
}

func GetGames() ([]types.Product, error) {
	var games []types.Product
	if err := dbGames.All(&games); err != nil {
		return nil, err
	}
	return games, nil
}

func GetAccounts() ([]types.Account, error) {
	var accounts []types.Account
	if err := dbAccounts.All(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// CreateAccountFromSignIn creates a new account from the sign in state
func CreateAccountFromSignIn(state types.GogAuthenticationChrome) error {
	account, err := NewAccountFromLogin(state)
	if err != nil {
		return err
	}

	return dbAccounts.Save(account)
}

// NewAccountFromLogin creates a new account from the login state
func NewAccountFromLogin(state types.GogAuthenticationChrome) (*types.Account, error) {
	account := types.Account{
		Cookies: []types.Cookie{},
	}

	if state.User.Email == "" {
		return nil, errors.New("email is empty")
	}

	if state.User.Username == "" {
		return nil, errors.New("username is empty")
	}

	if state.User.UserID == "" {
		return nil, errors.New("user id is empty")
	}

	if state.User.GalaxyUserID == "" {
		return nil, errors.New("galaxy user id is empty")
	}

	if state.User.Avatar == "" {
		return nil, errors.New("avatar is empty")
	}

	if state.User.Country == "" {
		return nil, errors.New("country is empty")
	}

	if len(state.Cookies) == 0 {
		return nil, errors.New("cookies are empty")
	}

	account.Email = state.User.Email
	account.Username = state.User.Username
	account.UserID = state.User.UserID
	account.GalaxyUserID = state.User.GalaxyUserID
	account.Avatar = state.User.Avatar
	account.Country = state.User.Country

	for _, cookie := range state.Cookies {
		account.Cookies = append(account.Cookies, cookie)
	}

	return &account, nil
}
