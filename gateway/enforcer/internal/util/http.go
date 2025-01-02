package util

import (
	"net/http"
)

// MakeGETRequest HTTP client for making GET requests
func MakeGETRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}