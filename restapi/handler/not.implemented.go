package handler

import (
	"net/http"

	"github.com/7phs/coding-challenge-auction/restapi/errCode"
	"github.com/gin-gonic/gin"
)

func NotImplemented(c *gin.Context) {
	resp := &RespError{}
	resp.AddError(errCode.ErrNotImplemented, "not implemented")

	c.JSON(http.StatusNotImplemented, resp)
}
