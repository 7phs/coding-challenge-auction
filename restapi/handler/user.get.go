package handler

import (
	"fmt"
	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/errCode"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type UserBid struct {
	ItemId  models.ItemKey `json:"item_id"`
	Bid     float64        `json:"bid"`
	Updated time.Time      `json:"updated"`
}

type UserBidList []*UserBid

func (o *UserBidList) Add(bid models.BidRecI) {
	*o = append(*o, &UserBid{
		ItemId:  bid.ItemId(),
		Bid:     bid.Bid(),
		Updated: time.Unix(0, bid.Updated()),
	})
}

func (o *UserBidList) Sort() {
	sort.Slice(*o, func(i, j int) bool {
		if (*o)[i].Updated.After((*o)[j].Updated) {
			return true
		} else if (*o)[i].Updated.Before((*o)[j].Updated) {
			return false
		}

		if (*o)[i].Bid > (*o)[j].Bid {
			return true
		} else if (*o)[i].Bid < (*o)[j].Bid {
			return false
		}

		return (*o)[i].ItemId > (*o)[j].ItemId
	})
}

type UserGetHandler struct {
	request struct {
		userId int // :userID
	}
	response struct {
		RespError
		Data struct {
			Id   models.UserKey `json:"id"`
			Name string         `json:"name"`
			Bids UserBidList    `json:"bids"`
		} `json:"data"`
	}
}

func (o *UserGetHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	var err error

	o.request.userId, err = strconv.Atoi(c.Param("userID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "user_id: "+err.Error())
	}

	return
}

func (o *UserGetHandler) Validate() (errList ErrorRecordList) {
	if o.request.userId == 0 {
		errList.AddError(errCode.ErrValidation, "user_id: empty")
	}

	return
}

func UserGet(c *gin.Context) {
	handler := UserGetHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("user/get: failed to bind - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}

	logPrefix := fmt.Sprintf("user/get: #%d", handler.request.userId)

	log.Info(logPrefix + ", handle")
	// VALIDATE
	if err := handler.Validate(); err != nil {
		logrus.Error(logPrefix+", failed to validate params - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}
	// GET
	user, err := models.Users.Get(models.UserKey(handler.request.userId))
	if err != nil {
		logrus.Error(logPrefix+", failed to get - ", err)

		handler.response.AddError(errCode.ErrUserProcessed, err)
		c.JSON(http.StatusNotFound, handler.response)
		return
	}
	// MARSHALING BIDS LIST
	bidsList := user.Bids()

	result := make(UserBidList, 0, len(bidsList))
	for _, b := range bidsList {
		result.Add(b)
	}
	result.Sort()
	// RESPONSE
	handler.response.Data.Id = user.Id()
	handler.response.Data.Name = user.Name()
	handler.response.Data.Bids = result
	c.JSON(http.StatusCreated, handler.response)
}
