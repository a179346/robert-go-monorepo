package dbhelper

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	_ "github.com/lib/pq"
)

func Open() (*sql.DB, error) {
	dbConfig := post_board_config.GetDBConfig()

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	return sql.Open("postgres", databaseURL)
}

func WaitFor(ctx context.Context, db *sql.DB) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			_, err := db.Query("SELECT 1")
			if err == nil {
				return
			}
			console.Warnf("connecting to database: %v", err)
			time.Sleep(2 * time.Second)
		}
	}
}
