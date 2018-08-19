package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/handler"
)

type User struct{}

func NewUser() *User {
	return &User{}
}

func (o *User) Add(name string) (models.UserKey, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return 0, errors.New("name: empty")
	} else if len(name) > limitLen {
		return 0, errors.New("name: greater than " + strconv.Itoa(limitLen))
	}

	params := url.Values{}
	params.Add("name", name)

	ur := Endpoint() + "/user"
	body := params.Encode()

	req, err := http.NewRequest("POST", ur, bytes.NewReader([]byte(body)))
	if err != nil {
		return 0, fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := handler.UserAddHandler{}

	err = Execute(req, &resp.Response)
	if err != nil {
		return 0, err
	}

	if len(resp.Response.Errors) > 0 {
		return 0, resp.Response.Errors
	}

	return resp.Response.Data.Id, nil

}

func (o *User) Update(id models.UserKey, name string) error {
	if id == 0 {
		return errors.New("id: empty")
	}

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return errors.New("name: empty")
	} else if len(name) > limitLen {
		return errors.New("name: greater than " + strconv.Itoa(limitLen))
	}

	params := url.Values{}
	params.Add("name", name)

	ur := Endpoint() + fmt.Sprintf("/user/%d", id)
	body := params.Encode()

	req, err := http.NewRequest("PUT", ur, bytes.NewReader([]byte(body)))
	if err != nil {
		return fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := handler.ItemAddHandler{}

	err = Execute(req, &resp.Response)
	if err != nil {
		return err
	}

	if len(resp.Response.Errors) > 0 {
		return resp.Response.Errors
	}

	return nil
}

func (o *User) Get(id models.UserKey) (*handler.UserGetResponse, error) {
	if id == 0 {
		return nil, errors.New("id: empty")
	}

	ur := Endpoint() + fmt.Sprintf("/user/%d", id)

	req, err := http.NewRequest("GET", ur, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}

	resp := handler.UserGetResponse{}

	err = Execute(req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Errors) > 0 {
		return nil, resp.Errors
	}

	return &resp, nil
}
