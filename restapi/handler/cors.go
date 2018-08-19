package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.teamc.io/worldskills/esatk/api/configuration"
)

func AllowCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

func Options(methods string) func(*gin.Context) {
	if configuration.Conf.Cors {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Methods", methods)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "accept, content-type")
			c.String(http.StatusOK, "ok")
		}
	} else {
		return func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		}
	}
}
