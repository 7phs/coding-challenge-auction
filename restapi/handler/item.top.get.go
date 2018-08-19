package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/errCode"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type ItemTopGetHandler struct {
	request struct {
		itemId int // :itemID
	}
	response struct {
		RespError
		Data struct {
			UserId models.UserKey `json:"user_id"`
			Bid    float64        `json:"bid"`
		} `json:"data"`
	}
}

func (o *ItemTopGetHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	var err error

	o.request.itemId, err = strconv.Atoi(c.Param("itemID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "item_id: "+err.Error())
	}

	return
}

func (o *ItemTopGetHandler) Validate() (errList ErrorRecordList) {
	if o.request.itemId == 0 {
		errList.AddError(errCode.ErrValidation, "item_id: empty")
	}

	return
}

func ItemTopGet(c *gin.Context) {
	handler := ItemTopGetHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("item/top/get: failed to bind - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}

	logPrefix := fmt.Sprintf("item/top/get: #%d", handler.request.itemId)

	log.Info(logPrefix + ", handle")
	// VALIDATE
	if err := handler.Validate(); err != nil {
		logrus.Error(logPrefix+", failed to validate params - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}
	// GET ITEM
	item, err := models.Items.Get(models.ItemKey(handler.request.itemId))
	if err != nil {
		logrus.Error(logPrefix+", failed to get - ", err)

		handler.response.AddError(errCode.ErrItemProcessed, err)
		c.JSON(http.StatusNotFound, handler.response)
		return
	}
	// GET TOP
	top, err := item.Top()
	if err != nil {
		logrus.Error(logPrefix+", failed to get the top bid - ", err)

		handler.response.AddError(errCode.ErrItemProcessed, err)
		c.JSON(http.StatusNotFound, handler.response)
		return
	}

	handler.response.Data.UserId = top.UserId()
	handler.response.Data.Bid = top.Bid()
	c.JSON(http.StatusCreated, handler.response)
}
