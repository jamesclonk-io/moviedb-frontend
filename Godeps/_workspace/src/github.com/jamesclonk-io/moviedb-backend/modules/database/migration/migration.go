package migration

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/mattes/migrate/migrate"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.GetLogger()
}

func RunMigrations(dbUri, dbType string) {
	errors, ok := migrate.UpSync(dbUri, fmt.Sprintf("./migrations/%s", dbType))
	if !ok {
		for _, err := range errors {
			log.Error(err)
		}
		log.Fatal("Could not migrate up database")
	}
}
