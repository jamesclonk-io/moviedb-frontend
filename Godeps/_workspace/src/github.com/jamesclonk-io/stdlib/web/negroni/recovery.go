package negroni

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
)

type Recovery struct {
	Logger     *logrus.Logger
	PrintStack bool
}

func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     logger.GetLogger(),
		PrintStack: true,
	}
}

func (r *Recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]

			stacktrace := r.Logger.WithFields(logrus.Fields{
				"error": err,
				"stack": string(stack),
			})
			stacktrace.Error("PANIC")

			if r.PrintStack {
				fmt.Fprintf(rw, "PANIC: %s\n%s", err, stack)
			}
		}
	}()
	next(rw, req)
}
