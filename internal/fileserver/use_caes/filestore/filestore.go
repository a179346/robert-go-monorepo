package filestore_use_case

import (
	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type FileStoreUseCase struct {
	fileStorePather fileStorePather
}

func New(config fileserver_config.StoreConfig) FileStoreUseCase {
	return FileStoreUseCase{
		fileStorePather: newFileStorePather(config.RootPath),
	}
}

func (filestore FileStoreUseCase) AppendHandler(router *roberthttp.Router) {
	router.Handle("/download", filestore.downloadHandler)
}
