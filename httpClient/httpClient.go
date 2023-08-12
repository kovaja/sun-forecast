package httpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var httpc = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := httpc.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)

	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Call to %s failed. Code: %d", url, r.StatusCode))
	}

	return nil
}
