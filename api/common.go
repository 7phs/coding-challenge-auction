package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/7phs/coding-challenge-auction/helpers"
)

const (
	defaultAddr = "http://localhost:8080"
	limitLen    = 4000
)

func Endpoint() string {
	return helpers.GetEnv("ADDR", defaultAddr)
}

func Execute(req *http.Request, response interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request for '%s': %s", req.URL, err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return fmt.Errorf("failed to execute request for '%s': status %s", req.URL, resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to execute request for '%s': reading body err - %s", req.URL, err)
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("failed to execute request for '%s': unmarshal body err - %s", req.URL, err)
	}

	return nil
}
