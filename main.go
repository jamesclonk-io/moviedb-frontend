package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/moviedb-backend/modules/moviedb"
	"github.com/jamesclonk-io/moviedb-frontend/modules/navbar"
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
	backendUrl = env.Get("JCIO_MOVIEDB_BACKEND", "http://moviedb-backend.jamesclonk.io")
}

type NavigationElement struct {
	Name     string              `json:"name"`
	Link     string              `json:"link,omitempty"`
	Icon     string              `json:"icon,omitempty"`
	Dropdown []NavigationElement `json:"dropdown,omitempty"`
}

func main() {
	// setup http handler
	n := setup()

	// start web server
	server := web.NewServer()
	server.Start(n)
}

func setup() *negroni.Negroni {
	backendClient = web.NewBackendClient()

	frontend := web.NewFrontend("jamesclonk.io - Movie Database")

	// setup routes
	frontend.NewRoute("/", movies)
	frontend.NewRoute("/movies", movies)

	frontend.NewRoute("/actors", actors)
	frontend.NewRoute("/directors", directors)
	frontend.NewRoute("/statistics", statistics)

	frontend.NewRoute("/error/{.*}", createError)

	// setup navbar
	frontend.SetNavigation(navbar.GetNavigation())

	n := negroni.Sbagliato()
	n.UseHandler(frontend.Router)

	return n
}

func getData(f func(string, string) *web.Page, urlPart string, req *http.Request) *web.Page {
	var query string
	if len(req.URL.RawQuery) > 0 {
		query = "?" + req.URL.RawQuery
	}
	query = fmt.Sprintf("%s%s", urlPart, query)

	response, err := backendClient.Get(fmt.Sprintf("%s%s", backendUrl, query))
	if err != nil {
		return web.Error("Error!", http.StatusInternalServerError, err)
	}
	return f(response, query)
}

func movies(w http.ResponseWriter, req *http.Request) *web.Page {
	return getData(func(response, query string) *web.Page {
		var data []moviedb.MovieListing
		if err := json.Unmarshal([]byte(response), &data); err != nil {
			return web.Error("Error!", http.StatusInternalServerError, err)
		}
		return &web.Page{
			ActiveLink: query,
			Content:    data,
			Template:   "index",
		}
	}, "/movies", req)
}

func actors(w http.ResponseWriter, req *http.Request) *web.Page {
	return getData(func(response, query string) *web.Page {
		var data []moviedb.Person
		if err := json.Unmarshal([]byte(response), &data); err != nil {
			return web.Error("Error!", http.StatusInternalServerError, err)
		}
		return &web.Page{
			ActiveLink: query,
			Content:    data,
			Template:   "people",
		}
	}, "/actors", req)
}

func directors(w http.ResponseWriter, req *http.Request) *web.Page {
	return getData(func(response, query string) *web.Page {
		var data []moviedb.Person
		if err := json.Unmarshal([]byte(response), &data); err != nil {
			return web.Error("Error!", http.StatusInternalServerError, err)
		}
		return &web.Page{
			ActiveLink: query,
			Content:    data,
			Template:   "people",
		}
	}, "/directors", req)
}

func statistics(w http.ResponseWriter, req *http.Request) *web.Page {
	return getData(func(response, query string) *web.Page {
		var data moviedb.Statistics
		if err := json.Unmarshal([]byte(response), &data); err != nil {
			return web.Error("Error!", http.StatusInternalServerError, err)
		}
		return &web.Page{
			ActiveLink: query,
			Content:    data,
			Template:   "statistics",
		}
	}, "/statistics", req)
}

func createError(w http.ResponseWriter, req *http.Request) *web.Page {
	return web.Error(
		"jamesclonk.io - Movie Database - Error",
		http.StatusInternalServerError,
		fmt.Errorf("Error!"),
	)
}
