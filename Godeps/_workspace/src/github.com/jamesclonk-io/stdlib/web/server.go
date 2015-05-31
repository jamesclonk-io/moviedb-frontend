package web

import (
	"net"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/tylerb/graceful"
)

type Server struct {
	Logger   *logrus.Logger
	port     string
	certFile string
	keyFile  string
}

func NewServer() *Server {
	return &Server{
		Logger: logger.GetLogger(),
	}
}

func (s *Server) Start(handler http.Handler) {
	s.port = env.Get("PORT", "3000")
	s.certFile = env.Get("JCIO_HTTP_CERT_FILE", "")
	s.keyFile = env.Get("JCIO_HTTP_KEY_FILE", "")

	address := ":" + s.port
	s.Logger.WithField("address", address).Info("Start HTTP Server")

	server := &graceful.Server{
		Server:  &http.Server{Addr: address, Handler: handler},
		Timeout: 15 * time.Second,
	}

	var err error
	if len(s.certFile) > 0 && len(s.keyFile) > 0 {
		err = server.ListenAndServeTLS(s.certFile, s.keyFile)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			s.Logger.Fatal(err)
		}
	}
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) CertFile() string {
	return s.certFile
}

func (s *Server) KeyFile() string {
	return s.keyFile
}
