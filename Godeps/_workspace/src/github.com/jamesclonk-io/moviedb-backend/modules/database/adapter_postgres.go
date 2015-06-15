package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresAdapter(uri string) *Adapter {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	return &Adapter{db}
}
