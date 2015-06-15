package database

import (
	"database/sql"
)

type Adapter struct {
	*sql.DB
}
