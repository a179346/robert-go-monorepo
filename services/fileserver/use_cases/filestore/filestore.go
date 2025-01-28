package filestore_use_case

import (
	fileserver_config "github.com/a179346/robert-go-monorepo/services/fileserver/config"
	"github.com/gohf-http/gohf/v6"
)

type FileStoreUseCase struct {
	fileStoreQueries  fileStoreQueries
	fileStoreCommands fileStoreCommands
}

func New() FileStoreUseCase {
	fileStorePather := newFileStorePather(fileserver_config.GetStorageConfig().StoreRootPath)

	return FileStoreUseCase{
		fileStoreQueries:  newFileStoreQueries(fileStorePather),
		fileStoreCommands: newFileStoreCommands(fileStorePather),
	}
}

func (filestore FileStoreUseCase) AppendHandler(router *gohf.Router) {
	router.GET("/download", filestore.downloadHandler)
	router.POST("/upload/{id}", filestore.uploadHandler)
}
