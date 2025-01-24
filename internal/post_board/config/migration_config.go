package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type MigrationConfig struct {
	FolderPath string
	Verbose    bool
	Up         bool
}

var migrationConfig MigrationConfig

func init() {
	migrationConfig.FolderPath = env_helper.GetString("MIGRATION_FOLDER_PATH", "internal/post_board/migrations")
	migrationConfig.Verbose = env_helper.GetBool("MIGRATION_VERBOSE", true)
	migrationConfig.Up = env_helper.GetBool("MIGRATION_UP", true)
}

func GetMigrationConfig() MigrationConfig {
	return migrationConfig
}
