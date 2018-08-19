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
	UserId   models.UserKey `json:"user_id"`
	UserName string         `json:"user_name"`
	Bid      string         `json:"bid"`
	bid      float64
	Updated  time.Time `json:"updated"`
}

type ItemBidList []*ItemBid

func (o *ItemBidList) Add(bid models.BidRecI) {
	var (
		b    = bid.Bid()
		uId  = bid.UserId()
		u, _ = models.Users.Get(uId)
	)

	*o = append(*o, &ItemBid{
		UserId:   uId,
		UserName: u.Name(),
		Bid:      models.FormatBid(b),
		bid:      b,
		Updated:  time.Unix(0, bid.Updated()),
	})
}

func (o *ItemBidList) Sort() {
	sort.Slice(*o, func(i, j int) bool {
		if (*o)[i].bid > (*o)[j].bid {
			return true
		} else if (*o)[i].bid < (*o)[j].bid {
			return false
		}

		return (*o)[i].Updated.After((*o)[j].Updated)
	})
}

type ItemGetResponse struct {
	RespError
	Data struct {
		Id    models.ItemKey `json:"id"`
		Title string         `json:"title"`
		Bids  ItemBidList    `json:"bids"`
	} `json:"data"`
}

type ItemGetHandler struct {
	request struct {
		itemId int // :itemID
	}
	response ItemGetResponse
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
