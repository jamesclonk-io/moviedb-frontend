package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/jamesclonk-io/stdlib/web/negroni"
)

var (
	log           *logrus.Logger
	backendClient *web.BackendClient
	backendUrl    string
)

func init() {
	log = logger.GetLogger()
	backendClient = web.NewBackendClient()
	backendUrl = env.Get("JCIO_MOVIEDB_BACKEND", "http://moviedb-backend.jamesclonk.io")
}

func main() {
	frontend := web.NewFrontend("jamesclonk.io - Movie Database")
	frontend.NewRoute("/", index)

	n := negroni.Sbagliato()
	n.UseHandler(frontend.Router)

	server := web.NewServer()
	server.Start(n)
}

func index(w http.ResponseWriter, req *http.Request) *web.Page {
	response, err := backendClient.Get(backendUrl)
	if err != nil {
		return web.Error("Error!", http.StatusInternalServerError, err)
	}

	return &web.Page{
		ActiveLink: "/",
		Content:    response,
		Template:   "index",
	}
}
