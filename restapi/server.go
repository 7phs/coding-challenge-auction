package restapi

import (
"context"
"net/http"
"time"

"github.com/gin-gonic/gin"
log "github.com/sirupsen/logrus"
	"github.com/7phs/coding-challenge-auction/config"
	"github.com/7phs/coding-challenge-auction/restapi/handler"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 5 * time.Second
)

type Server struct {
	http.Server
}

func Init(stage string) {
	switch stage {
	case config.StageTesting:
		gin.SetMode(gin.TestMode)
	case config.StageProduction:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func NewServer(conf *config.Config) *Server {
	return &Server{
		Server: http.Server{
			Addr:         conf.Addr,
			Handler:      handler.DefaultRouter(),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (o *Server) Run() *Server {
	go func() {
		log.Info("http: start listening ", o.Addr)
		if err := o.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start listening ", o.Addr, ": ", err)
		}
	}()

	return o
}

func (o *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info("http: shutdown")
	if err := o.Server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown HTTP server")
	}
}
