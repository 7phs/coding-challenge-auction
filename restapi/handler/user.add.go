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

type UserAddHandler struct {
	request struct {
		Name string `form:"name"`
	}
	Response struct {
		RespError
		Data struct {
			Id   models.UserKey `json:"id"`
			Name string         `json:"name"`
		} `json:"data"`
	}
}

func (o *UserAddHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	if err := c.ShouldBindWith(&o.request, binding.Default(c.Request.Method, c.ContentType())); err != nil {
		errList.AddError(errCode.ErrBinding, "title: "+err.Error())
	}

	return
}

func (o *UserAddHandler) Validate() (errList ErrorRecordList) {
	o.request.Name = strings.TrimSpace(o.request.Name)
	if l := len(o.request.Name); l == 0 {
		errList.AddError(errCode.ErrValidation, "name: empty")
	} else if l > limitTitleLen {
		errList.AddError(errCode.ErrValidation, "name: length greater than a limit "+strconv.Itoa(limitTitleLen))
	}

	return
}

func UserAdd(c *gin.Context) {
	handler := UserAddHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("user/add: failed to bind - ", err)

		handler.Response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.Response)
		return
	}

	logPrefix := fmt.Sprintf("user/add: name '%s'", handler.request.Name)

	log.Info(logPrefix + ", handle")
	// VALIDATE
	if err := handler.Validate(); err != nil {
		logrus.Error(logPrefix+", failed to validate params - ", err)

		handler.Response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.Response)
		return
	}
	// PUSH
	user := models.User(handler.request.Name)

	log.Info(logPrefix + ", add into model")
	err := models.Users.Add(user)
	if err != nil {
		logrus.Error(logPrefix+", failed to add into model - ", err)

		handler.Response.AddError(errCode.ErrUserProcessed, err)
		c.JSON(http.StatusBadRequest, handler.Response)
		return
	}
	// RESPONSE
	handler.Response.Data.Id = user.Id()
	handler.Response.Data.Name = user.Name()
	c.JSON(http.StatusCreated, handler.Response)
}
