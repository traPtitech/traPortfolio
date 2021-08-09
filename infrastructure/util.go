package infrastructure

import (
	"fmt"
	"net/http"
	"os"
)

var (
	isDevelopment = (os.Getenv("APP_EMV") == "development")
)

func apiGet(client *http.Client, endpoint string, path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", endpoint, path), nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
