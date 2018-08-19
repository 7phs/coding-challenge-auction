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

type ItemAddHandler struct {
	request struct {
		Title string `form:"title"`
	}
	Response struct {
		RespError
		Data struct {
			Id    models.ItemKey `json:"id"`
			Title string         `json:"title"`
		} `json:"data"`
	}
}

func (o *ItemAddHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	if err := c.ShouldBindWith(&o.request, binding.Default(c.Request.Method, c.ContentType())); err != nil {
		errList.AddError(errCode.ErrBinding, "title: "+err.Error())
	}

	return
}

func (o *ItemAddHandler) Validate() (errList ErrorRecordList) {
	o.request.Title = strings.TrimSpace(o.request.Title)
	if l := len(o.request.Title); l == 0 {
		errList.AddError(errCode.ErrValidation, "title: empty")
	} else if l > limitTitleLen {
		errList.AddError(errCode.ErrValidation, "title: length greater than a limit "+strconv.Itoa(limitTitleLen))
	}

	return
}

func ItemAdd(c *gin.Context) {
	handler := ItemAddHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("item/add: failed to bind - ", err)

		handler.Response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.Response)
		return
	}

	logPrefix := fmt.Sprintf("item/add: title '%s'", handler.request.Title)

	log.Info(logPrefix + ", handle")
	// VALIDATE
	if err := handler.Validate(); err != nil {
		logrus.Error(logPrefix+", failed to validate params - ", err)

		handler.Response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.Response)
		return
	}
	// PUSH
	item := models.Item(handler.request.Title)

	log.Info(logPrefix + ", add into model")
	err := models.Items.Add(item)
	if err != nil {
		logrus.Error(logPrefix+", failed to add into model - ", err)

		handler.Response.AddError(errCode.ErrItemProcessed, err)
		c.JSON(http.StatusBadRequest, handler.Response)
		return
	}
	// RESPONSE
	handler.Response.Data.Id = item.Id()
	handler.Response.Data.Title = item.Title()
	c.JSON(http.StatusCreated, handler.Response)
}
