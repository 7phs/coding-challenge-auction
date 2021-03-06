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

type UserUpdateHandler struct {
	request struct {
		userId int // :userID

		Name string `form:"name"`
	}
	response struct {
		RespError
		Data struct {
			Id   models.UserKey `json:"id"`
			Name string         `json:"name"`
		} `json:"data"`
	}
}

func (o *UserUpdateHandler) Bind(c *gin.Context) (errList ErrorRecordList) {
	var err error

	o.request.userId, err = strconv.Atoi(c.Param("userID"))
	if err != nil {
		errList.AddError(errCode.ErrBinding, "user_id: "+err.Error())
	}

	if err := c.ShouldBindWith(&o.request, binding.Default(c.Request.Method, c.ContentType())); err != nil {
		errList.AddError(errCode.ErrBinding, "name: "+err.Error())
	}

	return
}

func (o *UserUpdateHandler) Validate() (errList ErrorRecordList) {
	if o.request.userId == 0 {
		errList.AddError(errCode.ErrValidation, "user_id: empty")
	}

	o.request.Name = strings.TrimSpace(o.request.Name)
	if l := len(o.request.Name); l == 0 {
		errList.AddError(errCode.ErrValidation, "name: empty")
	} else if l > limitTitleLen {
		errList.AddError(errCode.ErrValidation, "name: length greater than a limit "+strconv.Itoa(limitTitleLen))
	}

	return
}

func UserUpdate(c *gin.Context) {
	handler := UserUpdateHandler{}
	// BIND
	if err := handler.Bind(c); err != nil {
		logrus.Error("user/update: failed to bind - ", err)

		handler.response.AppendError(err)
		c.JSON(http.StatusBadRequest, handler.response)
		return
	}

	logPrefix := fmt.Sprintf("user/update: #%d; name '%s'", handler.request.userId, handler.request.Name)

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
	// UPDATE
	log.Info(logPrefix + ", update a record")

	user.SetName(handler.request.Name)
	// RESPONSE
	handler.response.Data.Id = user.Id()
	handler.response.Data.Name = user.Name()
	c.JSON(http.StatusCreated, handler.response)
}
