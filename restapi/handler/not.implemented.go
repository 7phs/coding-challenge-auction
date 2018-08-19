package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.teamc.io/worldskills/esatk/api/errCode"
)

func NotImplemented(c *gin.Context) {
	resp := &RespError{}
	resp.AddError(errCode.ErrNorImplemented, "not implemented")

	c.JSON(http.StatusNotImplemented, resp)
}
