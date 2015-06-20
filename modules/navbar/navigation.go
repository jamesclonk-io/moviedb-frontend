package navbar

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/moviedb-backend/modules/moviedb"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

var (
	log   *logrus.Logger
	chars = []rune("abcdefghijklmnopqrstuvwxyz")
)

func init() {
	log = logger.GetLogger()
}

func GetNavigation() web.Navigation {
	moviesNav := web.Navigation{
		web.NavigationElement{
			Name: "by Name",
			Link: "/movies?sort=title&by=asc",
		},
		web.NavigationElement{
			Name: "by Score",
			Link: "/movies?sort=score&by=desc&sort=title&by=asc",
		},
		web.NavigationElement{
			Name: "by Rating",
			Link: "/movies?sort=rating&by=desc&sort=title&by=asc",
		},
		web.NavigationElement{
			Name: "by Year",
			Link: "/movies?sort=year&by=desc&sort=title&by=asc",
		},
		web.NavigationElement{
			Name: "Divider",
			Link: "#",
		},
	}
	for i := 5; i > 0; i-- {
		element := web.NavigationElement{
			Name: strings.Repeat("✰", i),
			Link: fmt.Sprintf("/movies?query=score&value=%d&sort=title&by=asc", i),
		}
		moviesNav = append(moviesNav, element)
	}

	titlesNav := web.Navigation{
		web.NavigationElement{
			Name: "0-9",
			Link: "/movies?query=char&value=num&sort=title&by=asc",
		},
	}
	for _, c := range chars {
		element := web.NavigationElement{
			Name: strings.ToUpper(string(c)),
			Link: fmt.Sprintf("/movies?query=char&value=%s&sort=title&by=asc", string(c)),
		}
		titlesNav = append(titlesNav, element)
	}

	return web.Navigation{
		web.NavigationElement{
			Name:     "Movies",
			Link:     "#",
			Icon:     "fa-film",
			Dropdown: moviesNav,
		},
		web.NavigationElement{
			Name:     "Titles",
			Link:     "#",
			Icon:     "fa-book",
			Dropdown: titlesNav,
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
}

func getGenres() ([]moviedb.Genre, error) {
	backendClient := web.NewBackendClient()
	backendUrl := env.Get("JCIO_MOVIEDB_BACKEND", "http://moviedb-backend.jamesclonk.io")

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
