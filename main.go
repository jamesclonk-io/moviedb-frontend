package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/moviedb-backend/modules/moviedb"
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

type NavigationElement struct {
	Name     string              `json:"name"`
	Link     string              `json:"link,omitempty"`
	Icon     string              `json:"icon,omitempty"`
	Dropdown []NavigationElement `json:"dropdown,omitempty"`
}

func main() {
	frontend := web.NewFrontend("jamesclonk.io - Movie Database")
	frontend.NewRoute("/", index)
	frontend.NewRoute("/movies", movies)

	// TODO: refactor this into own func !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	nav := web.Navigation{
		web.NavigationElement{
			Name: "Movies",
			Link: "#",
			Icon: "fa-film",
			Dropdown: web.Navigation{
				web.NavigationElement{
					Name: "by Name",
					Link: "/movies/by_name",
				},
				web.NavigationElement{
					Name: "by Score",
					Link: "/movies/by_score",
				},
				web.NavigationElement{
					Name: "by Rating",
					Link: "/movies/by_rating",
				},
				web.NavigationElement{
					Name: "by Year",
					Link: "/movies/by_year",
				},
				web.NavigationElement{
					Name: "Divider",
					Link: "#",
				},
				web.NavigationElement{
					Name: "✰✰✰✰✰",
					Link: "/movies/by_5stars",
				},
				web.NavigationElement{
					Name: "✰✰✰✰",
					Link: "/movies/by_4stars",
				},
				web.NavigationElement{
					Name: "✰✰✰",
					Link: "/movies/by_3stars",
				},
				web.NavigationElement{
					Name: "✰✰",
					Link: "/movies/by_2stars",
				},
				web.NavigationElement{
					Name: "✰",
					Link: "/movies/by_1stars",
				},
			},
		},
		web.NavigationElement{
			Name: "Titles",
			Link: "#",
			Icon: "fa-book",
			Dropdown: web.Navigation{
				web.NavigationElement{
					Name: "0-9",
					Link: "/movie_titles/by_num",
				},
				web.NavigationElement{
					Name: "A",
					Link: "/movie_titles/by_a",
				},
				web.NavigationElement{
					Name: "B",
					Link: "/movie_titles/by_b",
				},
			},
		},
		web.NavigationElement{
			Name:     "Genres",
			Link:     "#",
			Icon:     "fa-heartbeat",
			Dropdown: getGenreNavigation(),
		},
		web.NavigationElement{
			Name: "People",
			Link: "#",
			Icon: "fa-users",
			Dropdown: web.Navigation{
				web.NavigationElement{
					Name: "Actors",
					Link: "/actors",
				},
				web.NavigationElement{
					Name: "Directors",
					Link: "/directors",
				},
			},
		},
		web.NavigationElement{
			Name: "Statistics",
			Link: "/statistics",
			Icon: "fa-bar-chart",
		},
	}
	frontend.SetNavigation(nav)

	n := negroni.Sbagliato()
	n.UseHandler(frontend.Router)

	server := web.NewServer()
	server.Start(n)
}

func getGenreNavigation() web.Navigation {
	genres, err := getGenres()
	if err != nil {
		entry := log.WithFields(logrus.Fields{
			"error": err,
			"info":  "Could not get genres from backend",
		})
		entry.Error("Loading genres")
		return nil
	}

	var nav web.Navigation
	for _, genre := range genres {
		element := web.NavigationElement{
			Name: genre.Name,
			Link: fmt.Sprintf("/movies?query=genre&value=%d", genre.Id),
		}
		nav = append(nav, element)
	}
	return nav
}

func getGenres() ([]moviedb.Genre, error) {
	// TODO: refactor GET & UNMARSHAL into own func (withContext-style) !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	response, err := backendClient.Get(backendUrl + "/genres")
	if err != nil {
		return nil, err
	}

	var genres []moviedb.Genre
	if err := json.Unmarshal([]byte(response), &genres); err != nil {
		return nil, err
	}
	return genres, nil
}

func index(w http.ResponseWriter, req *http.Request) *web.Page {
	response, err := backendClient.Get(backendUrl + "/movies")
	if err != nil {
		return web.Error("Error!", http.StatusInternalServerError, err)
	}

	return &web.Page{
		ActiveLink: "/",
		Content:    response,
		Template:   "index",
	}
}

func movies(w http.ResponseWriter, req *http.Request) *web.Page {
	response, err := backendClient.Get(backendUrl + "/movies?" + req.URL.RawQuery)
	if err != nil {
		return web.Error("Error!", http.StatusInternalServerError, err)
	}

	return &web.Page{
		ActiveLink: "/",
		Content:    response,
		Template:   "index",
	}
}
