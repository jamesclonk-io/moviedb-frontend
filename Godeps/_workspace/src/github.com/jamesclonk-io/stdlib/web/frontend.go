package web

import (
	j "encoding/json"
	"html/template"
	"net/http"

	"github.com/unrolled/render"
)

// Frontend has control over routing, rendering and page mastership.
type Frontend struct {
	Title      string
	Router     *Router
	Render     *render.Render
	PageMaster *PageMaster
}

func NewFrontend(title string) *Frontend {
	r := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Extensions: []string{".html"},
		Funcs: []template.FuncMap{template.FuncMap{
			"isEvenNumber": isEvenNumber,
			"isOddNumber":  isOddNumber,
			"html":         htmlize,
			"json":         jsonize,
		}},
	})
	router := NewRouter()
	pm := &PageMaster{title, nil, nil, "index", http.StatusOK, nil}

	f := &Frontend{title, router, r, pm}

	// default 404
	router.NotFoundHandler = f.NotFoundHandler(title)

	return f
}

func (f *Frontend) SetNavigation(nav Navigation) {
	f.PageMaster.Navigation = nav
}

func isEvenNumber(input int) bool {
	return input%2 == 0
}

func isOddNumber(input int) bool {
	return !isEvenNumber(input)
}

func htmlize(input string) template.HTML {
	return template.HTML(input)
}

func jsonize(input interface{}) template.JS {
	bytes, err := j.Marshal(input)
	if err != nil {
		return ""
	}
	return template.JS(bytes)
}
