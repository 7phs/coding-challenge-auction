package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi/errCode"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type ItemUpdateHandler struct {
	request struct {
		itemId int // :itemID

		Title string `form:"title"`
	}
	response struct {
		RespError
		Data struct {
			Id    models.ItemKey `json:"id"`
			Title string         `json:"title"`
		} `json:"data"`
	}
}

func (o *ItemUpdateHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	var err error

	o.request.itemId, err = strconv.Atoi(c.Param("itemID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "item_id: "+err.Error())
	}

	if err := c.ShouldBindWith(&o.request, binding.Default(c.Request.Method, c.ContentType())); err != nil {
		errList.AddError(errCode.ErrBinding, "title: "+err.Error())
	}

	return
}

func (o *ItemUpdateHandler) Validate() (errList ErrorRecordList) {
	if o.request.itemId == 0 {
		errList.AddError(errCode.ErrValidation, "item_id: empty")
	}

	o.request.Title = strings.TrimSpace(o.request.Title)
	if l := len(o.request.Title); l == 0 {
		errList.AddError(errCode.ErrValidation, "title: empty")
	} else if l > limitTitleLen {
		errList.AddError(errCode.ErrValidation, "title: length greater than a limit "+strconv.Itoa(limitTitleLen))
	}

	return
}

func ItemUpdate(c *gin.Context) {
	handler := ItemUpdateHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("item/update: failed to bind - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}

	logPrefix := fmt.Sprintf("item/update: #%d; title '%s'", handler.request.itemId, handler.request.Title)

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
	// UPDATE
	log.Info(logPrefix + ", update a record")

	item.SetTitle(handler.request.Title)
	// RESPONSE
	handler.response.Data.Id = item.Id()
	handler.response.Data.Title = item.Title()
	c.JSON(http.StatusCreated, handler.response)
}
