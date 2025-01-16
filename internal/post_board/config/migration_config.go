package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type MigrationConfig struct {
	FolderPath string
	Verbose    bool
	Up         bool
}

var migrationConfig MigrationConfig

func initMigrationConfig() {
	migrationConfig.FolderPath = env_helper.GetStringEnv("MIGRATION_FOLDER_PATH", "internal/post_board/database/migrations")
	migrationConfig.Verbose = env_helper.GetBoolEnv("MIGRATION_VERBOSE", true)
	migrationConfig.Up = env_helper.GetBoolEnv("MIGRATION_UP", true)
}

func GetMigrationConfig() MigrationConfig {
	initAll()
	return migrationConfig
}
