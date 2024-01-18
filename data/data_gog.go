package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/davidroman0O/gogog/types"
)

type GogClient struct {
	*http.Client
}

func NewGogClient() *GogClient {
	j, _ := cookiejar.New(nil)
	return &GogClient{
		Client: &http.Client{
			Jar:       j,
			Timeout:   time.Minute * 5,
			Transport: types.NewTransporter(generateUserAgent()),
		},
	}
}

func SetCookies(client *http.Client, cks []*types.Cookie, siteUrl string) error {
	var cookies []*http.Cookie
	for _, cookie := range cks {
		name := cookie.Name
		value := cookie.Value
		cookie := &http.Cookie{
			Domain: cookie.Domain,
			Name:   name,
			Path:   cookie.Path,
			Secure: cookie.Secure,
			Value:  value,
		}
		cookies = append(cookies, cookie)
	}

	u, err := url.Parse(siteUrl)
	if err != nil {
		return err
	}

	client.Jar.SetCookies(u, cookies)
	return nil
}

func CheckCookies(client *http.Client, siteUrl string) (bool, error) {
	req, err := client.Get(siteUrl + "/userData.json")
	if err != nil {
		return false, err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return false, errors.New(req.Status)
	}
	var obj types.UserData
	err = json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return false, err
	}
	ok := obj.IsLoggedIn
	if ok {
		fmt.Println("Signed in as " + obj.Username + ".\n")
	}
	return ok, nil
}
