package negroni

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	classico "github.com/codegangsta/negroni"
	"github.com/jamesclonk-io/stdlib/logger"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{logger.GetLogger()}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()

	entry := l.WithFields(logrus.Fields{
		"request": req.URL.RequestURI(),
		"method":  req.Method,
		"remote":  req.RemoteAddr,
	})
	if id := req.Header.Get("X-Request-Id"); len(id) > 0 {
		entry = entry.WithField("request_id", id)
	}
	if id := req.Header.Get("X-Cf-Requestid"); len(id) > 0 {
		entry = entry.WithField("cf_request_id", id)
	}
	if id := req.Header.Get("X-Vcap-Request-Id"); len(id) > 0 {
		entry = entry.WithField("vcap_request_id", id)
	}
	entry.Info("Handling request")

	next(rw, req)

	duration := time.Since(start)
	res := rw.(classico.ResponseWriter)
	entry.WithFields(logrus.Fields{
		"status":      res.Status(),
		"text_status": http.StatusText(res.Status()),
		"duration":    duration,
	}).Info("Completed request")
}
