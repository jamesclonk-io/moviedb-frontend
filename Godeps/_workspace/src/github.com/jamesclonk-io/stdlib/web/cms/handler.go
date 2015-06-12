package cms

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/jamesclonk-io/stdlib/web"
)

func (c *CMS) ViewHandler(w http.ResponseWriter, req *http.Request) *web.Page {
	filename := path.Base(req.URL.RequestURI())

	// find file
	var html template.HTML
	for _, c := range c.Data.Content {
		if path.Join("/", c.Path, c.Basename) == req.URL.RequestURI() {
			html = c.Content
		}
	}

	// return 404 if content was not found
	if html == "" {
		return &web.Page{
			Title:      fmt.Sprintf("%s - %s - Not Found", c.Title, filename),
			StatusCode: http.StatusNotFound,
			Template:   "404",
		}
	}

	// wrap into struct
	content := struct {
		Title string
		HTML  template.HTML
	}{
		Title: filename,
		HTML:  html,
	}

	return &web.Page{
		Title:      fmt.Sprintf("%s - %s", c.Title, filename),
		ActiveLink: req.URL.RequestURI(),
		Content:    content,
		Template:   "cms",
	}
}

func (c *CMS) RefreshHandler(w http.ResponseWriter, req *http.Request) *web.Page {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.checkData(true); err != nil {
		return web.Error(c.Title, http.StatusInternalServerError, err)
	}

	return &web.Page{
		Title:      c.Title,
		ActiveLink: "/",
		Content:    nil,
		Template:   "index",
	}
}
