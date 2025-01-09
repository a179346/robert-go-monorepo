package filestore_use_case

import (
	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
)

type FileStoreUseCase struct {
	fileStorePather fileStorePather
}

func New(config fileserver_config.StorageConfig) FileStoreUseCase {
	return FileStoreUseCase{
		fileStorePather: newFileStorePather(config.StoreRootPath),
	}
}

func (filestore FileStoreUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("GET /download", filestore.downloadHandler)
	router.Handle("POST /upload/{id}", filestore.uploadHandler)
}
