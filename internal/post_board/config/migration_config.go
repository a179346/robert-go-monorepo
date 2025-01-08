package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type MigrationConfig struct {
	FolderPath string
	Verbose    bool
	Up         bool
}

func newMigrationConfig() MigrationConfig {
	folderPath := env_helper.GetStringEnv("MIGRATION_FOLDER_PATH", "internal/post_board/database/migrations")
	verbose := env_helper.GetBoolEnv("MIGRATION_VERBOSE", true)
	up := env_helper.GetBoolEnv("MIGRATION_UP", true)

	return MigrationConfig{
		FolderPath: folderPath,
		Verbose:    verbose,
		Up:         up,
	}
}
