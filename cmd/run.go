package cmd

import (
	"os"
	"os/signal"

	"github.com/7phs/coding-challenge-auction/config"
	"github.com/7phs/coding-challenge-auction/logger"
	"github.com/7phs/coding-challenge-auction/models"
	"github.com/7phs/coding-challenge-auction/restapi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()

		log.Info(ApplicationInfo())

		config.Init()

		models.Init()

		restapi.Init(config.Stage)

		// Create new server instance
		srv := restapi.
			NewServer(&config.Conf).
			Run()

		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)
		<-stop
		log.Info("interrupt signal")

		srv.Shutdown()
		models.Shutdown()

		log.Info("finish")
	},
}
