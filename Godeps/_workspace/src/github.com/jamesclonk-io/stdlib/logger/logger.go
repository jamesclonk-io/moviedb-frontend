package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
)

func init() {
	checkEnv()
}

func checkEnv() {
	if env.Get("JCIO_ENV", "") == "production" || // manual
		env.Get("VCAP_APPLICATION", "") != "" || // cf / lattice
		env.Get("DYNO", "") != "" { // heroku
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
