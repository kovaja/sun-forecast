package httpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var httpc = &http.Client{Timeout: 10 * time.Second}

func readJsonResponse(url string, r *http.Response, target interface{}) error {
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Call to %s failed. Code: %d", url, r.StatusCode))
	}

	err := json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func GetJsonWithAuth(url string, token string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	r, err := httpc.Do(req)
	if err != nil {
		return err
	}

	return readJsonResponse(url, r, target)
}

func GetJson(url string, target interface{}) error {
	r, err := httpc.Get(url)
	if err != nil {
		return err
	}

	return readJsonResponse(url, r, target)
}
