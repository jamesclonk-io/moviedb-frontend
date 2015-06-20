package web

import (
	"fmt"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) *Page

func (f *Frontend) NewHandler(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		page := fn(w, req)

		if len(page.Title) == 0 {
			page.Title = f.PageMaster.Title
		}
		if page.Content == nil {
			page.Content = f.PageMaster.Content
		}
		if page.Data == nil {
			page.Data = f.PageMaster.Data
		}
		if len(page.Template) == 0 {
			page.Template = f.PageMaster.Template
		}
		if page.StatusCode == 0 {
			page.StatusCode = f.PageMaster.StatusCode
		}
		if page.Navigation == nil {
			page.Navigation = f.PageMaster.Navigation
		}

		if page.Error != nil {
			f.Render.HTML(w, page.StatusCode, "error", page)
			return
		}
		f.Render.HTML(w, page.StatusCode, page.Template, page)
	}
}

func (b *Backend) NewHandler(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		page := fn(w, req)

		if page.StatusCode == 0 {
			page.StatusCode = http.StatusOK
		}

		if page.Headers != nil {
			for key, values := range page.Headers {
				w.Header().Del(key)
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
		}

		if page.Error != nil {
			b.Render.JSON(w, page.StatusCode, page.Error.Error())
			return
		}
		b.Render.JSON(w, page.StatusCode, page.Content)
	}
}

func (b *Backend) NewSecuredHandler(fn Handler) http.HandlerFunc {
	f := b.NewHandler(fn)
	return func(w http.ResponseWriter, req *http.Request) {
		// do auth
		if b.hmac {
			// hmac check when using HMAC Backend
			ok, err := b.hmacAuth(req)
			if err != nil {
				b.Render.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
				return
			}
			if !ok {
				b.Render.JSON(w, http.StatusUnauthorized, "Unauthorized!")
				return
			}
		} else {
			// basic auth check when using TLS Backend
			user, password, ok := req.BasicAuth()
			if !ok || user != b.user || password != b.password {
				b.Render.JSON(w, http.StatusUnauthorized, "Unauthorized!")
				return
			}
		}
		f(w, req)
	}
}

func (f *Frontend) NotFoundHandler(title string) http.HandlerFunc {
	return f.NewHandler(func(http.ResponseWriter, *http.Request) *Page {
		return &Page{
			Title:      title,
			StatusCode: http.StatusNotFound,
			Template:   "404",
		}
	})
}

func (b *Backend) NotFoundHandler(headers http.Header) http.HandlerFunc {
	return b.NewHandler(func(http.ResponseWriter, *http.Request) *Page {
		return &Page{
			Headers:    headers,
			StatusCode: http.StatusNotFound,
			Content:    map[string]string{"Code": "404", "Status": "This is not the JSON you are looking for.."},
		}
	})
}

func Error(title string, code int, err error) *Page {
	return &Page{
		Title:      title,
		StatusCode: code,
		Error:      err,
		Template:   "error",
	}
}
