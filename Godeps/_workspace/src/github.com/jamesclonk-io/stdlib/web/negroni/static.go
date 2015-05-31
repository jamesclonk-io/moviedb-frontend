package negroni

import "net/http"

type Static struct {
	Dir http.FileSystem
}

func NewStatic() *Static {
	return &Static{
		Dir: http.Dir("public"),
	}
}

func (s *Static) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}
	file := r.URL.Path

	f, err := s.Dir.Open(file)
	if err != nil {
		next(rw, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		next(rw, r)
		return
	}

	if fi.IsDir() {
		next(rw, r)
		return
	}

	// cache all static content for 1 week
	rw.Header().Set("Cache-control", "public, max-age=604800")

	http.ServeContent(rw, r, file, fi.ModTime(), f)
}
