package filestore_use_case

import (
	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
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
	router.Handle("GET /download", filestore.downloadHandler)
	router.Handle("POST /upload/{id}", filestore.uploadHandler)
}
