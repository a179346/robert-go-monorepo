package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

type MigrationConfig struct {
	FolderPath string
	Verbose    bool
	Up         bool
}

var migrationConfig MigrationConfig

func init() {
	migrationConfig.FolderPath = envhelper.GetString("MIGRATION_FOLDER_PATH", "services/post_board/migrations")
	migrationConfig.Verbose = envhelper.GetBool("MIGRATION_VERBOSE", true)
	migrationConfig.Up = envhelper.GetBool("MIGRATION_UP", true)
}

func GetMigrationConfig() MigrationConfig {
	return migrationConfig
}
