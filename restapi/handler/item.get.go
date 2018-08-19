package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/errCode"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type ItemBid struct {
	UserId  models.UserKey `json:"user_id"`
	Bid     float64        `json:"bid"`
	Updated time.Time      `json:"updated"`
}

type ItemBidList []*ItemBid

func (o *ItemBidList) Add(bid models.BidRecI) {
	*o = append(*o, &ItemBid{
		UserId:  bid.UserId(),
		Bid:     bid.Bid(),
		Updated: time.Unix(0, bid.Updated()),
	})
}

func (o *ItemBidList) Sort() {
	sort.Slice(*o, func(i, j int) bool {
		if (*o)[i].Bid > (*o)[j].Bid {
			return true
		} else if (*o)[i].Bid < (*o)[j].Bid {
			return false
		}

		return (*o)[i].Updated.After((*o)[j].Updated)
	})
}

type ItemGetHandler struct {
	request struct {
		itemId int // :itemID
	}
	response struct {
		RespError
		Data struct {
			Id    models.ItemKey `json:"id"`
			Title string         `json:"title"`
			Bids  []*ItemBid     `json:"bids"`
		} `json:"data"`
	}
}

func (o *ItemGetHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	var err error

	o.request.itemId, err = strconv.Atoi(c.Param("itemID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "item_id: "+err.Error())
	}

	return
}

func (o *ItemGetHandler) Validate() (errList ErrorRecordList) {
	if o.request.itemId == 0 {
		errList.AddError(errCode.ErrValidation, "item_id: empty")
	}

	return
}

func ItemGet(c *gin.Context) {
	handler := ItemGetHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("item/get: failed to bind - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}

	logPrefix := fmt.Sprintf("item/get: #%d", handler.request.itemId)

	log.Info(logPrefix + ", handle")
	// VALIDATE
	if err := handler.Validate(); err != nil {
		logrus.Error(logPrefix+", failed to validate params - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}
	// GET
	item, err := models.Items.Get(models.ItemKey(handler.request.itemId))
	if err != nil {
		logrus.Error(logPrefix+", failed to get - ", err)

		handler.response.AddError(errCode.ErrItemProcessed, err)
		c.JSON(http.StatusNotFound, handler.response)
		return
	}
	// MARSHALING BIDS LIST
	bidsList := item.Bids()

	result := make(ItemBidList, 0, len(bidsList))
	for _, b := range bidsList {
		result.Add(b)
	}
	result.Sort()
	// RESPONSE
	handler.response.Data.Id = item.Id()
	handler.response.Data.Title = item.Title()
	handler.response.Data.Bids = result
	c.JSON(http.StatusCreated, handler.response)
}
