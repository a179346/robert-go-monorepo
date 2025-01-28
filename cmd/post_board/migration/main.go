package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	"github.com/a179346/robert-go-monorepo/services/post_board/database/dbhelper"
	_ "github.com/a179346/robert-go-monorepo/services/post_board/migrations"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		console.Errorf("%s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := dbhelper.Open()
	if err != nil {
		return fmt.Errorf("opendb.Open error: %w", err)
	}
	defer db.Close()

	dbhelper.WaitFor(ctx, db)
	if ctx.Err() != nil {
		return ctx.Err()
	}

	migrationConfig := post_board_config.GetMigrationConfig()

	gooseProvider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS(migrationConfig.FolderPath),
		goose.WithVerbose(true),
		goose.WithAllowOutofOrder(true),
	)
	if err != nil {
		return fmt.Errorf("goose.NewProvider error: %w", err)
	}

	if migrationConfig.Up {
		if _, err := gooseProvider.Up(ctx); err != nil {
			return fmt.Errorf("gooseProvider.Up error: %w", err)
		}
	} else {
		if _, err := gooseProvider.Down(ctx); err != nil {
			return fmt.Errorf("gooseProvider.Down error: %w", err)
		}
	}

	return ctx.Err()
}
