package migrations

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up1736330008, Down1736330008)
}

func Up1736330008(ctx context.Context, tx *sql.Tx) error {
	// password: "secret123"
	query := []string{
		`INSERT INTO "user" ("id", "email", "name", "encrypted_pass")`,
		`VALUES ('8ff6fe28-14c5-4dc8-a0bf-749fa8a212a0', 'admin@google.com', 'admin', 'fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4');`,
	}

	_, err := tx.ExecContext(ctx, strings.Join(query, "\n"))
	if err != nil {
		return err
	}
	return nil
}

func Down1736330008(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM "user" WHERE "id" = '8ff6fe28-14c5-4dc8-a0bf-749fa8a212a0';`)
	if err != nil {
		return err
	}
	return nil
}
