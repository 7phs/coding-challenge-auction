package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/7phs/coding-challenge-auction/helpers"
	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/handler"
)

type Bid struct {
	itemId models.ItemKey
	userId models.UserKey
}

func NewBid(itemId models.ItemKey, userId models.UserKey) *Bid {
	return &Bid{
		itemId: itemId,
		userId: userId,
	}
}

func (o *Bid) Validated() error {
	var errList helpers.ErrList

	if o.itemId == 0 {
		errList.Add(errors.New("item_id: empty"))
	}

	if o.userId == 0 {
		errList.Add(errors.New("user_id: empty"))
	}

	if len(errList) == 0 {
		return nil
	}

	return errList
}

func (o *Bid) Push(bid float64) error {
	if err := o.Validated(); err != nil {
		return err
	}

	if bid <= 0 {
		return errors.New("bid: should be positive and great than zero")
	}

	params := url.Values{}
	params.Add("bid", models.FormatBid(bid))

	ur := Endpoint() + fmt.Sprintf("/bid/%d/user/%d", o.itemId, o.userId)
	body := params.Encode()

	req, err := http.NewRequest("POST", ur, bytes.NewReader([]byte(body)))
	if err != nil {
		return fmt.Errorf("failed to create a request for '%s': %s", ur, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := handler.BidPushHandler{}

	err = Execute(req, &resp.Response)
	if err != nil {
		return err
	}

	if len(resp.Response.Errors) > 0 {
		return resp.Response.Errors
	}

	return nil
}
