package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
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

func SetCookies(client *http.Client, cks []types.Cookie, siteUrl string) error {
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

func FetchUser(client *http.Client, siteUrl string) (*types.UserData, bool, error) {
	req, err := client.Get(siteUrl + "/userData.json")
	if err != nil {
		return nil, false, err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return nil, false, errors.New(req.Status)
	}
	var obj types.UserData
	err = json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return nil, false, err
	}
	ok := obj.IsLoggedIn
	if ok {
		fmt.Println("Signed in as " + obj.Username + ".\n")
	}
	return &obj, ok, nil
}

func GetGameMeta(client *http.Client, siteUrl string, id int) (*types.GameMeta, error) {
	req, err := client.Get(
		fmt.Sprintf("%v/account/gameDetails/%v.json", siteUrl, strconv.Itoa(id)),
	)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return nil, errors.New(req.Status)
	}

	var obj types.GameMeta
	err = json.NewDecoder(req.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func Search(client *http.Client, siteUrl string, search types.SearchParams) ([]types.Product, error) {
	var err error
	var req *http.Request
	if req, err = http.NewRequest(
		http.MethodGet, fmt.Sprintf("%v/account/getFilteredProducts", siteUrl), nil); err != nil {
		return nil, err
	}

	pageNum := 1
	query := url.Values{}
	// i don't know what's it's used for
	query.Set("hiddenFlag", "0")

	if search.Language != nil {
		if !strings.Contains(*search.Language, "all") {
			query.Set("language", *search.Language)
		}
	}

	query.Set("mediaType", "1")
	// query.Set("mediaType", "game")
	// query.Set("sortBy", "date_purchased")
	query.Set("sortBy", "title")

	if search.Query != nil {
		query.Set("search", *search.Query)
	}

	if search.PlatformName != nil {
		query.Set("system", string(*search.PlatformName))
	}

	// query.Set("totalPages", "1")

	var products []types.Product

	for {
		query.Set("page", strconv.Itoa(pageNum))
		req.URL.RawQuery = query.Encode()

		do, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if do.StatusCode != http.StatusOK {
			do.Body.Close()
			return nil, errors.New(do.Status)
		}

		var obj types.Search
		err = json.NewDecoder(do.Body).Decode(&obj)
		do.Body.Close()
		if err != nil {
			return nil, err
		}

		if obj.TotalPages == 0 {
			break
		}

		products = append(products, obj.Products...)
		if pageNum == obj.TotalPages {
			break
		}

		pageNum++
		time.Sleep(time.Second * 1)
	}
	return products, nil
}

// func DownloadImage(client *http.Client, siteUrl string, product types.Product) error {

// 	url := fmt.Sprintf("https:%v_product_card_v2_logo_480x285.png", product.Image)
// 	destination := fmt.Sprintf("./tmp/imgs/%v/cover.png", product.ID)
// 	// if err := CheckAndCreateDataFolder(fmt.Sprintf("imgs/%v", product.ID)); err != nil {
// 	// 	return err
// 	// }

// 	fmt.Println(url)

// 	// Send an HTTP GET request to the image URL
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Create the destination file
// 	file, err := os.Create(destination)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Copy the response body to the file
// 	_, err = io.Copy(file, resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
