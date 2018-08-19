package handler

import (
	"net/http"

	"github.com/7phs/coding-challenge-auction/config"
	"github.com/gin-gonic/gin"
)

func AllowCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

func Options(methods string) func(*gin.Context) {
	if config.Conf.Cors {
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
