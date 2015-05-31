package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Env_Get(t *testing.T) {
	assert.NotEqual(t, "foobar", Get("PATH", "foobar"))
	assert.Equal(t, "bar", Get("foo", "bar"))
}

func Test_Env_MustGet(t *testing.T) {
	path := MustGet("PATH")
	assert.NotNil(t, path)

	called := false
	defer func() {
		err := recover()
		if err != nil {
			assert.Equal(t, "Required env var [evil_foo_bar] is missing!", err)
			called = true
		}
	}()
	MustGet("evil_foo_bar")
	assert.True(t, called)
}
