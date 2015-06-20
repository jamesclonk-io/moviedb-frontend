package database

import (
	"database/sql"

	"github.com/JamesClonk/vcap"
	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/moviedb-backend/modules/database/migration"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
)

var (
	log *logrus.Logger
)

type Adapter struct {
	Database *sql.DB
	URI      string
	Type     string
}

func init() {
	log = logger.GetLogger()
}

func NewAdapter() (db *Adapter) {
	var databaseType, databaseUri string

	// get db type
	databaseType = env.Get("JCIO_DATABASE_TYPE", "postgres")

	// check for VCAP_SERVICES first
	data, err := vcap.New()
	if err != nil {
		log.Fatal(err)
	}
	if service := data.GetService("moviedb"); service != nil {
		if uri, ok := service.Credentials["uri"]; ok {
			databaseUri = uri.(string)
		}
	}

	// if JCIO_DATABASE_URL is not yet set then try to read it from ENV
	if len(databaseUri) == 0 {
		databaseUri = env.MustGet("JCIO_DATABASE_URI")
	}

	// setup database adapter
	switch databaseType {
	case "postgres":
		db = newPostgresAdapter(databaseUri)
	case "sqlite":
		db = newSQLiteAdapter(databaseUri)
	default:
		log.Fatalf("Invalid database type: %s\n", databaseType)
	}

	// panic if no database adapter was set up
	if db == nil {
		log.Fatal("Could not set up database adapter")
	}

	// run db migrations
	migration.RunMigrations(databaseUri, databaseType)

	return db
}
