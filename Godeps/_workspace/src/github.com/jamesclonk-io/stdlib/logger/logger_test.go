package logger

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Logger_GetLogger(t *testing.T) {
	assert.Equal(t, &logrus.TextFormatter{}, GetLogger().Formatter)
}

func Test_Logger_init(t *testing.T) {
	os.Setenv("JCIO_ENV", "production")
	checkEnv()
	assert.Equal(t, &logrus.JSONFormatter{}, GetLogger().Formatter)
}
