package negroni

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	classico "github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
)

func Test_Negroni_Logger(t *testing.T) {
	buf := bytes.NewBufferString("")
	rec := httptest.NewRecorder()

	l := NewLogger()
	l.Logger.Out = buf
	n := classico.New()
	n.Use(l)
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	req.Header.Add("X-Request-Id", "abc123")
	req.Header.Add("X-Cf-Requestid", "def 456")
	req.Header.Add("X-Vcap-Request-Id", "ghi 789")
	assert.Nil(t, err)
	n.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.NotEqual(t, 0, len(buf.String()))
	assert.True(t, strings.Contains(buf.String(), `msg="Handling request"`))
	assert.True(t, strings.Contains(buf.String(), `msg="Completed request"`))
	assert.True(t, strings.Contains(buf.String(), `text_status="Not Found"`))
	assert.True(t, strings.Contains(buf.String(), `request_id=abc123`))
	assert.True(t, strings.Contains(buf.String(), `cf_request_id="def 456"`))
	assert.True(t, strings.Contains(buf.String(), `vcap_request_id="ghi 789"`))
}
