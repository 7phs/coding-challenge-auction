package logger

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/7phs/coding-challenge-auction/helpers"
)

var (
	logLevelMap = map[string]logrus.Level{
		"debug":   logrus.DebugLevel,
		"info":    logrus.InfoLevel,
		"warning": logrus.WarnLevel,
		"error":   logrus.ErrorLevel,
	}
)

type logWriter struct{}

func (o *logWriter) Write(p []byte) (n int, err error) {
	p = bytes.Replace(p, []byte{'\n'}, []byte{' '}, -1)

	logrus.Print(string(p))
	return len(p), nil
}

func Init() {
	// gin specific
	log.SetOutput(&logWriter{})
	gin.DefaultWriter = &logWriter{}
	gin.DefaultErrorWriter = &logWriter{}
	// init logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	// set a log level
	logLevelStr := strings.ToLower(helpers.GetEnv("LOG_LEVEL", "debug"))
	logLevel, ok := logLevelMap[logLevelStr]
	if !ok {
		logLevel = logLevelMap["error"]
	}
	logrus.SetLevel(logLevel)
}
