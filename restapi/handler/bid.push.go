package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/errCode"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type BidPushHandler struct {
	request struct {
		itemId int // :itemID
		userId int // :userID

		Bid float64 `form:"bid"`
	}
	response struct {
		RespError
	}
}

func (o *BidPushHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	var err error

	o.request.itemId, err = strconv.Atoi(c.Param("itemID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "item_id: "+err.Error())
	}

	o.request.userId, err = strconv.Atoi(c.Param("userID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "user_id: "+err.Error())
	}

	if err := c.ShouldBindWith(&o.request, binding.Default(c.Request.Method, c.ContentType())); err != nil {
		errList.AddError(errCode.ErrBinding, "bid: "+err.Error())
	}

	return
}

func (o *BidPushHandler) Validate() (errList ErrorRecordList) {
	if o.request.itemId == 0 {
		errList.AddError(errCode.ErrValidation, "item_id: empty")
	}

	if o.request.userId == 0. {
		errList.AddError(errCode.ErrValidation, "user_id: empty")
	}

	if o.request.Bid == 0. {
		errList.AddError(errCode.ErrValidation, "bid: equal 0")
	}

	return
}

func BidPush(c *gin.Context) {
	handler := BidPushHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("bid/push: failed to bind - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}

	logPrefix := fmt.Sprint("bid/push: item #", handler.request.itemId,
		"; user #", handler.request.userId,
		"; bid ", strconv.FormatFloat(handler.request.Bid, 'f', int(models.Precision), 64))

	log.Info(logPrefix + ", handle")
	// VALIDATE
	if err := handler.Validate(); err != nil {
		logrus.Error(logPrefix+", failed to validate params - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}
	// PUSH
	bidRec, created := models.Bids.Push(
		models.ItemKey(handler.request.itemId),
		models.UserKey(handler.request.userId),
		handler.request.Bid)
	log.Info(logPrefix+", model push ", bidRec)
	// RESPONSE
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	c.JSON(status, handler.response)
}
