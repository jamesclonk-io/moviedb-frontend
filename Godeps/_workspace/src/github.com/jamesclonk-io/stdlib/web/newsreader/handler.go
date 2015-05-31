package newsreader

import (
	"fmt"
	"net/http"

	"github.com/jamesclonk-io/stdlib/web"
)

func (n *NewsReader) ViewHandler(w http.ResponseWriter, req *http.Request) *web.Page {
	feeds, err := n.GetFeeds()
	if err != nil {
		return web.Error(fmt.Sprintf("%s - News Error", n.Title), http.StatusInternalServerError, err)
	}

	return &web.Page{
		Title:      fmt.Sprintf("%s - News", n.Title),
		ActiveLink: req.URL.RequestURI(),
		Content:    feeds,
		Template:   "news",
	}
}
