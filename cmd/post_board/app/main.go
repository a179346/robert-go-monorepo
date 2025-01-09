package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/opendb"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
)

func main() {
	config := post_board_config.New()

	db, err := opendb.Open(config.DB)
	if err != nil {
		log.Println("Error occurred: sql.Open")
		log.Fatal(err)
	}
	defer db.Close()

	userProvider := user_provider.New(db)

	data, err := userProvider.FindByEmail(
		context.Background(),
		"admin@gmail.com",
	)
	if err != nil {
		log.Println("Error occurred: userProvider")
		log.Fatal(err)
	}

	jsonText, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonText))
}
