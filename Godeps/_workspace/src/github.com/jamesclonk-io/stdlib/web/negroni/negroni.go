package negroni

import (
	classico "github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
)

type Negroni struct {
	*classico.Negroni
}

func Classico() *Negroni {
	n := classico.New()
	return &Negroni{n}
}

func Sbagliato() *Negroni {
	n := Classico()
	n.Use(NewRecovery())
	n.Use(NewLogger())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(NewStatic())
	return n
}

func (n *Negroni) Mescolare() *Negroni {
	n.Use(NewRecovery())
	n.Use(NewLogger())
	n.Use(NewStatic())
	return n
}
