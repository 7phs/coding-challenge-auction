package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/7phs/coding-challenge-auction/config"
)

func DefaultRouter() http.Handler {
	log.Info("restapi/router: init")

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	if config.Conf.Cors {
		router.Use(AllowCors())
	}
	// BID
	router.POST("/bid", NotImplemented)
	router.OPTIONS("/bid", Options("POST"))
	router.PUT("/bid/:GUID", NotImplemented)
	router.OPTIONS("/bid/:GUID", Options("PUT"))
	// ITEMS
	router.POST("/item", NotImplemented)
	router.OPTIONS("/item", Options("POST"))
	router.GET("/item/:GUID", NotImplemented)
	router.GET("/item/:GUID/top", NotImplemented)
	router.PUT("/item/:GUID", NotImplemented)
	router.OPTIONS("/item/:GUID", Options("PUT"))
	// USERS
	router.POST("/user", NotImplemented)
	router.OPTIONS("/user", Options("POST"))
	router.GET("/user/:GUID", NotImplemented)
	router.PUT("/user/:GUID", NotImplemented)
	router.OPTIONS("/user/:GUID", Options("PUT"))
	// HEALTH CHECK
	router.GET("/health/check", HealthCheck)

	return router
}
