package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/davidroman0O/gogog/types"
)

var (
	ErrUnreachable = errors.New("api is unreachable")
)

func Ping() error {
	var err error
	var resp *http.Response
	if resp, err = http.Get("http://localhost:8080/api/v1/ping"); err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		responseBody := string(body)
		if responseBody != "pong" {
			return errors.Join(ErrUnreachable, fmt.Errorf("answer is not pong"))
		}
		return nil
	}
	return errors.Join(ErrUnreachable, fmt.Errorf("couldn't ping the api"))
}

func GetAccounts() ([]types.Account, error) {
	var accounts []types.Account

	if err := dbAccounts.All(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func PostAccount(state types.GogAuthenticationChrome) error {
	var err error
	var resp *http.Response

	jsonData, err := json.Marshal(state)
	if err != nil {
		return err
	}

	bodyBuffer := bytes.NewBuffer(jsonData)

	if resp, err = http.Post("http://localhost:8080/api/v1/accounts", "application/json", bodyBuffer); err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		return nil
	}
	return errors.New("couldn't post the account")
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
		account.Cookies = append(account.Cookies, *cookie)
	}

	return &account, nil
}
