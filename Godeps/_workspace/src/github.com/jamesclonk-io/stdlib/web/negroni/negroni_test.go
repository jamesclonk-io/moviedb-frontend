package negroni

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Negroni_Classico(t *testing.T) {
	n := Classico()
	assert.Equal(t, 0, len(n.Handlers()))
}

func Test_Negroni_Sbagliato(t *testing.T) {
	n := Sbagliato()
	assert.Equal(t, 4, len(n.Handlers()))
	assert.Equal(t, "*negroni.Recovery", reflect.TypeOf(n.Handlers()[0]).String())
	assert.Equal(t, "*negroni.Logger", reflect.TypeOf(n.Handlers()[1]).String())
	assert.Equal(t, "*gzip.handler", reflect.TypeOf(n.Handlers()[2]).String())
	assert.Equal(t, "*negroni.Static", reflect.TypeOf(n.Handlers()[3]).String())
}

func Test_Negroni_Mescolare(t *testing.T) {
	n := Classico()
	n.Mescolare()
	assert.Equal(t, 3, len(n.Handlers()))
	assert.Equal(t, "*negroni.Recovery", reflect.TypeOf(n.Handlers()[0]).String())
	assert.Equal(t, "*negroni.Logger", reflect.TypeOf(n.Handlers()[1]).String())
	assert.Equal(t, "*negroni.Static", reflect.TypeOf(n.Handlers()[2]).String())
}
