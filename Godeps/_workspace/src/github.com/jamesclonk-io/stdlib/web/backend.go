package web

import (
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/unrolled/render"
)

// Backend has control over routing and rendering
type Backend struct {
	Router         *Router
	Render         *render.Render
	user, password string
	secret         []byte
	hmac           bool
}

func NewBackend() *Backend {
	var b *Backend
	// use HMAC based Backend if secret is set, otherwise assume TLS
	if len(env.Get("JCIO_HTTP_HMAC_SECRET", "")) != 0 {
		b = newHMACBackend()
	} else {
		b = newTLSBackend()
	}

	b.Render = render.New(render.Options{
		IndentJSON: false,
	})
	b.Router = NewRouter()
	b.Router.NotFoundHandler = b.NotFoundHandler(nil)

	return b
}

func newTLSBackend() *Backend {
	// enforce TLS certs / HTTPS listener for backend
	env.MustGet("JCIO_HTTP_CERT_FILE")
	env.MustGet("JCIO_HTTP_KEY_FILE")

	// backend needs user & password for basic auth
	user := env.MustGet("JCIO_HTTP_AUTH_USER")
	password := env.MustGet("JCIO_HTTP_AUTH_PASSWORD")

	return &Backend{
		user:     user,
		password: password,
		hmac:     false,
	}
}

func newHMACBackend() *Backend {
	// need preshared secret for HMAC backends
	secret := env.MustGet("JCIO_HTTP_HMAC_SECRET")

	return &Backend{
		secret: []byte(secret),
		hmac:   true,
	}
}
