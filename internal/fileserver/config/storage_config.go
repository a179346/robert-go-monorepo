package fileserver_config

import (
	"path/filepath"

	"github.com/a179346/robert-go-monorepo/pkg/env_helper"
)

type StorageConfig struct {
	isInited      bool
	RootPath      string
	StoreRootPath string
}

var storageConfig StorageConfig

func initStorageConfig() {
	if storageConfig.isInited {
		return
	}
	storageConfig.RootPath = env_helper.GetStringEnv("STORAGE_ROOT_PATH", "./storage/fileserver")

	storageConfig.StoreRootPath = filepath.Join(storageConfig.RootPath, "store")

	storageConfig.isInited = true
}

func GetStorageConfig() StorageConfig {
	initStorageConfig()
	return storageConfig
}
