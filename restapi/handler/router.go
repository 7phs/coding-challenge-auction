package handler

import (
	"net/http"

	"github.com/7phs/coding-challenge-auction/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	limitTitleLen = 4000
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
	router.POST("/bid/:itemID/user/:userID", BidPush)
	router.OPTIONS("/bid/:itemID/user/:userID", Options("POST"))
	// ITEMS
	router.POST("/item", ItemAdd)
	router.OPTIONS("/item", Options("POST"))
	router.GET("/item/:itemID", ItemGet)
	router.GET("/item/:itemID/top", ItemTopGet)
	router.PUT("/item/:itemID", ItemUpdate)
	router.OPTIONS("/item/:itemID", Options("PUT"))
	// USERS
	router.POST("/user", UserAdd)
	router.OPTIONS("/user", Options("POST"))
	router.GET("/user/:userID", UserGet)
	router.PUT("/user/:userID", UserUpdate)
	router.OPTIONS("/user/:userID", Options("PUT"))
	// HEALTH CHECK
	router.GET("/health/check", HealthCheck)

	return router
}
