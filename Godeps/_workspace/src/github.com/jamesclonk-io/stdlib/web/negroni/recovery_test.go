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

func Test_Negroni_Recovery(t *testing.T) {
	buf := bytes.NewBufferString("")
	rec := httptest.NewRecorder()

	n := classico.New()
	r := NewRecovery()
	r.Logger.Out = buf
	n.Use(r)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("!")
	}))

	n.ServeHTTP(rec, (*http.Request)(nil))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotEqual(t, 0, rec.Body.Len())
	assert.NotEqual(t, 0, len(buf.String()))
	assert.True(t, strings.Contains(buf.String(), `error="!"`))
}
