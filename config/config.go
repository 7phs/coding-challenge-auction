package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/7phs/coding-challenge-auction/helpers"
	"strings"
)

const (
	StageProduction  = "production"
	StageDevelopment = "development"
	StageTesting     = "testing"
)

const (
	defaultAddr  = ":8080"
	defaultCors  = "false"
	defaultStage = StageDevelopment
)

var (
	Conf  Config
	Stage string
)

type Config struct {
	Addr     string
	Cors     bool
	LogLevel string
}

func Init() {
	log.Info("config: init")

	Conf.Addr = helpers.GetEnv("ADDR", defaultAddr)
	Conf.Cors = strings.ToLower(helpers.GetEnv("CORS", defaultCors)) == "true"
	Stage = strings.ToLower(helpers.GetEnv("STAGE", defaultStage))

	validate()
}

func validate() {
	switch Stage {
	case StageProduction, StageDevelopment, StageTesting:
	default:
		log.Fatal("unsupported stage '" + Stage + "'. Supported: [" + strings.Join([]string{
			StageProduction,
			StageDevelopment,
			StageTesting,
		}, ", ") + "]")
	}
}
