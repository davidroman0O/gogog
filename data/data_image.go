package data

import (
	"encoding/base64"
	"io"
	"net/http"
)

func FetchImageToBase64(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	return base64Image, nil
}
