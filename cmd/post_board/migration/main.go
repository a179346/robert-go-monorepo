package main

import (
	"log"

	_ "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/dbhelper"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrationConfig := post_board_config.GetMigrationConfig()
	sourceURL := "file://" + migrationConfig.FolderPath

	db, err := dbhelper.Open()
	if err != nil {
		log.Fatalf("opendb.Open error: %v", err)
	}
	defer db.Close()
	dbhelper.WaitFor(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("postgres.WithInstance error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, post_board_config.GetDBConfig().Database, driver)
	if err != nil {
		log.Fatalf("migrate.NewWithDatabaseInstance error: %v", err)
	}
	m.Log = NewMigationLogger(migrationConfig.Verbose)

	if migrationConfig.Up {
		err = m.Up()
		if err != nil {
			if err.Error() == "no change" {
				log.Println("No change")
				return
			}

			log.Fatalf("m.Up error: %v", err)
		}
	} else {
		err = m.Steps(-1)
		if err != nil {
			log.Fatalf("m.Steps(-1) error: %v", err)
		}
	}
}

type MigationLogger struct {
	verbose bool
}

func NewMigationLogger(verbose bool) MigationLogger {
	return MigationLogger{verbose}
}

func (logger MigationLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (logger MigationLogger) Verbose() bool {
	return logger.verbose
}
