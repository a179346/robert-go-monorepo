package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/dbhelper"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	migrationConfig := post_board_config.GetMigrationConfig()
	sourceURL := "file://" + migrationConfig.FolderPath

	db, err := dbhelper.Open()
	if err != nil {
		return fmt.Errorf("opendb.Open error: %w", err)
	}
	defer db.Close()

	dbhelper.WaitFor(ctx, db)
	if ctx.Err() != nil {
		return ctx.Err()
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance error: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, post_board_config.GetDBConfig().Database, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance error: %w", err)
	}
	m.Log = NewMigationLogger(migrationConfig.Verbose)

	go func() {
		<-ctx.Done()
		log.Println("Gracefully shutting down ...")
		m.GracefulStop <- true
	}()

	if migrationConfig.Up {
		err = m.Up()
		if err != nil && err.Error() != "no change" {
			return fmt.Errorf("m.Up error: %w", err)
		}
	} else {
		err = m.Steps(-1)
		if err != nil {
			return fmt.Errorf("m.Steps(-1) error: %w", err)
		}
	}

	return ctx.Err()
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
