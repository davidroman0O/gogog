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
