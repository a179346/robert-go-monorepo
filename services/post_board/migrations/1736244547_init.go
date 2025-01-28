package migrations

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up1736244547, Down1736244547)
}

func Up1736244547(ctx context.Context, tx *sql.Tx) error {
	query := []string{
		`CREATE TABLE "user"`,
		`(`,
		`		"id" UUID NOT NULL,`,
		`		"email" VARCHAR(255) NOT NULL CONSTRAINT "user_email_key" UNIQUE,`,
		`		"name" VARCHAR(127) NOT NULL,`,
		`		"encrypted_pass" VARCHAR(255) NOT NULL,`,
		`		"created_at" TIMESTAMP NOT NULL DEFAULT NOW(),`,
		`		"updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),`,
		`		PRIMARY KEY ("id")`,
		`);`,

		`CREATE TABLE "post"`,
		`(`,
		`		"id" UUID NOT NULL,`,
		`		"author_id" UUID NOT NULL REFERENCES "user"("id") ON DELETE CASCADE,`,
		`		"content" TEXT NOT NULL,`,
		`		"created_at" TIMESTAMP NOT NULL DEFAULT NOW(),`,
		`		"updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),`,
		`		PRIMARY KEY ("id")`,
		`);`,
	}

	_, err := tx.ExecContext(ctx, strings.Join(query, "\n"))
	if err != nil {
		return err
	}
	return nil
}

func Down1736244547(ctx context.Context, tx *sql.Tx) error {
	query := []string{
		`DROP TABLE "post";`,
		`DROP TABLE "user";`,
	}

	_, err := tx.ExecContext(ctx, strings.Join(query, "\n"))
	if err != nil {
		return err
	}
	return nil
}
