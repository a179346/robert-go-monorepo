package dbhelper

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	_ "github.com/lib/pq"
)

func Open(dbConfig post_board_config.DBConfig) (*sql.DB, error) {
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

func WaitFor(db *sql.DB) {
	for {
		_, err := db.Query("SELECT 1")
		if err == nil {
			return
		}
		log.Printf("connecting to database: %v", err)
		time.Sleep(2 * time.Second)
	}
}
