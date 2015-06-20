package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web/negroni"
	"github.com/stretchr/testify/assert"
)

var (
	m *negroni.Negroni
)

func init() {
	os.Setenv("PORT", "3008")
	logrus.SetOutput(ioutil.Discard)
	logger.GetLogger().Out = ioutil.Discard

	os.Setenv("JCIO_MOVIEDB_BACKEND", "http://moviedb-backend.jamesclonk.io")
	os.Setenv("JCIO_HTTP_HMAC_SECRET", "who cares?")

	m = setup()
}

func Test_Main_404(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/something", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusNotFound, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database</title>`)
	assert.Contains(t, body, `<div class="alert alert-warning">This is not the page you are looking for..</div>`)
}

func Test_Main_500(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/error/something", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database - Error</title>`)
	assert.Contains(t, body, `<div class="alert alert-danger">Error: Error!</div>`)
}

func Test_Main_Index(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database</title>`)
	assert.Contains(t, body, `<li class=''><a href="/movies?query=score&amp;value=5&amp;sort=title&amp;by=asc">✰✰✰✰✰</a></li>`)
	assert.Contains(t, body, `<li class=''><a href="/movies?query=char&amp;value=h&amp;sort=title&amp;by=asc">H</a></li>`)
	assert.Contains(t, body, `<li class=''><a href="/movies?query=genre&amp;value=23">Crime</a></li>`)
	assert.Contains(t, body, `<td style="width:5%"><a class="no-underline" href="/movies?query=year&value=2010"><span class="label label-primary">2010</span></a></td>`)
	assert.Contains(t, body, `<td style="width:8%"><a class="no-underline score" href="/movies?query=score&value=4&sort=title&by=asc"><strong>★★★★</strong></a></td>`)
	assert.Contains(t, body, `<td><a class="no-underline" href="/movie/1026">Army of Darkness</a></td>`)
}

func Test_Main_MovieSort(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/movies?sort=year&by=asc&sort=title&by=asc", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<thead>
      <tr>
        <th>Year</th>
        <th>Score</th>
        <th>Title</th>
      </tr>
    </thead>
    <tbody>
      
      <tr>
        <td style="width:5%"><a class="no-underline" href="/movies?query=year&value=1962"><span class="label label-primary">1962</span></a></td>
        <td style="width:8%"><a class="no-underline score" href="/movies?query=score&value=4&sort=title&by=asc"><strong>★★★★</strong></a></td>
        <td><a class="no-underline" href="/movie/130">James Bond 007: Dr. No</a></td>
      </tr>`)
}

func Test_Main_OneStarMovies(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/movies?query=score&value=1&sort=title&by=asc", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.NotContains(t, body, `★★★★★`)
	assert.NotContains(t, body, `★★★★`)
	assert.NotContains(t, body, `★★★`)
	assert.NotContains(t, body, `★★`)
	assert.Contains(t, body, `<td><a class="no-underline" href="/movie/451">Eragon</a></td>`)
}

func Test_Main_TitleZMovies(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/movies?query=char&value=z&sort=title&by=asc", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.NotContains(t, body, `Argo`)
	assert.Contains(t, body, `Zatôichi`)
}

func Test_Main_GenreMusicMovies(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/movies?query=genre&value=32", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.NotContains(t, body, `Terminator`)
	assert.Contains(t, body, `O Brother, Where Art Thou?`)
}
