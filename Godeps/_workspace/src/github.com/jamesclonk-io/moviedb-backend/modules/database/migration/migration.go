package migration

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/moviedb-backend/modules/database"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/mattes/migrate/migrate"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.GetLogger()
}

func RunMigrations(basePath string, adapter *database.Adapter) {
	errors, ok := migrate.UpSync(adapter.URI, fmt.Sprintf("%s/%s", basePath, adapter.Type))
	if !ok {
		for _, err := range errors {
			log.Error(err)
		}
		log.Fatal("Could not migrate up database")
	}
}
