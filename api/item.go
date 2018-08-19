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

type Item struct{}

func NewItem() *Item {
	return &Item{}
}

func (o *Item) Add(title string) (models.ItemKey, error) {
	title = strings.TrimSpace(title)
	if len(title) == 0 {
		return 0, errors.New("title: empty")
	} else if len(title) > limitLen {
		return 0, errors.New("title: greater than " + strconv.Itoa(limitLen))
	}

	params := url.Values{}
	params.Add("title", title)

	ur := Endpoint() + "/item"
	body := params.Encode()

	req, err := http.NewRequest("POST", ur, bytes.NewReader([]byte(body)))
	if err != nil {
		return 0, fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := handler.ItemAddHandler{}

	err = Execute(req, &resp.Response)
	if err != nil {
		return 0, err
	}

	if len(resp.Response.Errors) > 0 {
		return 0, resp.Response.Errors
	}

	return resp.Response.Data.Id, nil

}

func (o *Item) Update(id models.ItemKey, title string) error {
	if id == 0 {
		return errors.New("id: empty")
	}

	title = strings.TrimSpace(title)
	if len(title) == 0 {
		return errors.New("title: empty")
	} else if len(title) > limitLen {
		return errors.New("title: greater than " + strconv.Itoa(limitLen))
	}

	params := url.Values{}
	params.Add("title", title)

	ur := Endpoint() + fmt.Sprintf("/item/%d", id)
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

func (o *Item) Get(id models.ItemKey) (*handler.ItemGetResponse, error) {
	if id == 0 {
		return nil, errors.New("id: empty")
	}

	ur := Endpoint() + fmt.Sprintf("/item/%d", id)

	req, err := http.NewRequest("GET", ur, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}

	resp := handler.ItemGetResponse{}

	err = Execute(req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Errors) > 0 {
		return nil, resp.Errors
	}

	return &resp, nil
}

func (o *Item) Top(id models.ItemKey) (*handler.ItemTopResponse, error) {
	if id == 0 {
		return nil, errors.New("id: empty")
	}

	ur := Endpoint() + fmt.Sprintf("/item/%d/top", id)

	req, err := http.NewRequest("GET", ur, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}

	resp := handler.ItemTopResponse{}

	err = Execute(req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Errors) > 0 {
		return nil, resp.Errors
	}

	return &resp, nil
}
