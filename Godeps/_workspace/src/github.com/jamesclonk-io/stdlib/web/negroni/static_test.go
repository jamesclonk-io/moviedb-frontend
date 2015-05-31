package negroni

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Negroni_Static_GET(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Body = new(bytes.Buffer)

	n := Classico()
	s := NewStatic()
	s.Dir = http.Dir(".")
	n.Use(s)

	req, err := http.NewRequest("GET", "http://localhost:3000/_fixtures/testfile", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "", rec.Header().Get("Expires"))
	assert.Equal(t, "public, max-age=604800", rec.Header().Get("Cache-control"))
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
	assert.Equal(t, "test data :)\n", rec.Body.String())
}

func Test_Negroni_Static_HEAD(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Body = new(bytes.Buffer)

	n := Classico()
	s := NewStatic()
	s.Dir = http.Dir(".")
	n.Use(s)
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("HEAD", "http://localhost:3000/_fixtures/testfile", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "", rec.Header().Get("Expires"))
	assert.Equal(t, "public, max-age=604800", rec.Header().Get("Cache-control"))
	if rec.Body.Len() != 0 {
		t.Errorf("Got non-empty body for HEAD request")
	}
}

func Test_Negroni_Static_POST(t *testing.T) {
	rec := httptest.NewRecorder()

	n := Classico()
	s := NewStatic()
	s.Dir = http.Dir(".")
	n.Use(s)
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("POST", "http://localhost:3000/_fixtures/testfile", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for POST request")
	}
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}

func Test_Negroni_Static_NoFile(t *testing.T) {
	rec := httptest.NewRecorder()

	n := Sbagliato()
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3000/_fixtures/testfile", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}

func Test_Negroni_Static_Dir(t *testing.T) {
	rec := httptest.NewRecorder()

	n := Classico()
	s := NewStatic()
	s.Dir = http.Dir(".")
	n.Use(s)
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3000/_fixtures", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}

func Test_Negroni_Static_DirWithSlash(t *testing.T) {
	rec := httptest.NewRecorder()

	n := Classico()
	s := NewStatic()
	s.Dir = http.Dir(".")
	n.Use(s)
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3000/_fixtures/", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}

func Test_Negroni_Static_BadDir(t *testing.T) {
	rec := httptest.NewRecorder()

	n := Classico()
	s := NewStatic()
	s.Dir = http.Dir("foobar")
	n.Use(s)
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3000/_fixtures/testfile", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}

func Test_Negroni_Static_WrongDir(t *testing.T) {
	rec := httptest.NewRecorder()

	n := Sbagliato()
	n.UseHandler(http.NotFoundHandler())

	req, err := http.NewRequest("GET", "http://localhost:3000/_fixtures/testfile", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
	if rec.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
	assert.Equal(t, "404 page not found\n", rec.Body.String())
}
