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
	assert.Contains(t, body, `<td style="width:5%"><a class="no-underline" href="/movies?query=year&value=2010"><span class="label label-default">2010</span></a></td>`)
	assert.Contains(t, body, `<td style="width:4%"><a class="no-underline" href="/movies?query=rating&value=16"><span class="label label-warning">16</span></a></td>`)
	assert.Contains(t, body, `<td style="width:5%"><a class="no-underline score" href="/movies?query=score&value=4&sort=title&by=asc"><strong>★★★★</strong></a></td>`)
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
      <th>Rating</th>
      <th>Score</th>
      <th>Title</th>
    </tr>
  </thead>
  <tbody>
    
    <tr>
      <td style="width:5%"><a class="no-underline" href="/movies?query=year&value=1962"><span class="label label-default">1962</span></a></td>
      <td style="width:4%"><a class="no-underline" href="/movies?query=rating&value=16"><span class="label label-warning">16</span></a></td>
      <td style="width:5%"><a class="no-underline score" href="/movies?query=score&value=4&sort=title&by=asc"><strong>★★★★</strong></a></td>
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

func Test_Main_Movie(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/movie/511", nil)
	if err != nil {
		t.Error(err)
	}
	req.RequestURI = "/movie/511"

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database - Apocalypse Now</title>`)
	assert.Contains(t, body, `<h3 class="list-group-item-heading">Apocalypse Now</h3>`)
	assert.Contains(t, body, `<p class="list-group-item-text">Apocalypse Now Redux</p>`)
	assert.Contains(t, body, `<img src="/images/movies/apocalypse_now.jpg" alt="Apocalypse Now"`)
	assert.Contains(t, body, `<a class="no-underline score" href="/movies?query=score&value=4&sort=title&by=asc"><strong>★★★★</strong></a>`)
	assert.Contains(t, body, `<a class="no-underline" href="/movies?query=genre&value=6&sort=title&by=asc">Drama</a>, <a class="no-underline" href="/movies?query=genre&value=15&sort=title&by=asc">War</a>`)
	assert.Contains(t, body, `<a class="no-underline" href="/movies?query=language&value=1&sort=title&by=asc">Deutsch</a>, <a class="no-underline" href="/movies?query=language&value=2&sort=title&by=asc">Englisch</a>`)
	assert.Contains(t, body, `It is the height of the war in Vietnam, and U.S. Army Captain Willard is sent by Colonel Lucas and a General to carry out a mission that, officially, &#039;does not exist - nor will it ever exist&#039;.`)
	assert.Contains(t, body, `<td><a class="no-underline" href="/person/1221">Albert Hall</a>, <a class="no-underline" href="/person/1224">Bo Byers</a>, <a class="no-underline" href="/person/1075">Dennis Hopper</a>, <a class="no-underline" href="/person/1219">Frederic Forrest</a>, <a class="no-underline" href="/person/1222">G.D. Spradlin</a>, <a class="no-underline" href="/person/489">Harrison Ford</a>, <a class="no-underline" href="/person/1225">James Keane</a>, <a class="no-underline" href="/person/1223">Jerry Ziesmer</a>, <a class="no-underline" href="/person/1226">Kerry Rossall</a>, <a class="no-underline" href="/person/50">Laurence Fishburne</a>, <a class="no-underline" href="/person/1217">Marlon Brando</a>, <a class="no-underline" href="/person/852">Martin Sheen</a>, <a class="no-underline" href="/person/1218">Robert Duvall</a>, <a class="no-underline" href="/person/1220">Sam Bottoms</a>, <a class="no-underline" href="/person/892">Scott Glenn</a>, </td>`)
	assert.Contains(t, body, `<td><a class="no-underline" href="/person/1216">Francis Ford Coppola</a>, </td>`)
}

func Test_Main_Person(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/person/211", nil)
	if err != nil {
		t.Error(err)
	}
	req.RequestURI = "/person/211"

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database - Quentin Tarantino</title>`)
	assert.Contains(t, body, `<h3 style="margin-bottom: 20px;">Quentin Tarantino</h3>`)
	assert.Contains(t, body, `<a href="#" class="list-group-item active"><h4 class="list-group-item-heading">Actor in:</h4></a>`)
	assert.Contains(t, body, `<td><a class="no-underline" href="/movie/51">From Dusk Till Dawn</a></td>`)
	assert.Contains(t, body, `<a href="#" class="list-group-item active"><h4 class="list-group-item-heading">Director of:</h4></a>`)
	assert.Contains(t, body, `<td><a class="no-underline" href="/movie/98">Kill Bill Vol.1</a></td>`)
}

func Test_Main_Actors(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/actors", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database</title>`)
	assert.Contains(t, body, `<div class="col-md-3 col-sm-4"><a class="no-underline" href="/person/171">Robert De Niro</a></div>
<div class="col-md-3 col-sm-4"><a class="no-underline" href="/person/778">Robert Downey Jr.</a></div>
<div class="col-md-3 col-sm-4"><a class="no-underline" href="/person/1218">Robert Duvall</a></div>`)
}

func Test_Main_Directors(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/directors", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database</title>`)
	assert.Contains(t, body, `<div class="col-md-3 col-sm-4"><a class="no-underline" href="/person/507">Yuen Wo Ping</a></div>
<div class="col-md-3 col-sm-4"><a class="no-underline" href="/person/512">Yu Wang</a></div>
<div class="col-md-3 col-sm-4"><a class="no-underline" href="/person/535">Zack Snyder</a></div>`)
}

func Test_Main_Statistics(t *testing.T) {
	response := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3008/statistics", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	assert.Contains(t, body, `<title>jamesclonk.io - Movie Database - Statistics</title>`)
	assert.Contains(t, body, `<td># Average Movies per day</td>`)
	assert.Contains(t, body, `<div id="top5actorsanddirectors"></div>`)
	assert.Contains(t, body, `var result = '<a href="/movies?query=score&value='+score+'&sort=title&by=asc" class="score no-underline">';`)
	assert.Contains(t, body, `\x22id\x22:396,\x22name\x22:\x22Clint Eastwood\x22`)
	assert.Contains(t, body, `\x22id\x22:493,\x22name\x22:\x22Kenji Kamiyama\x22`)
	assert.Contains(t, body, `\x22id\x22:483,\x22name\x22:\x22Bud Spencer\x22`)
	assert.Contains(t, body, `\x22count\x22`)
}
