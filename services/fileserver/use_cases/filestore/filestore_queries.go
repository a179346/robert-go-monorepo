package filestore_use_case

import (
	"errors"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/ztrue/tracerr"
)

type fileStoreQueries struct {
	fileStorePather fileStorePather
}

func newFileStoreQueries(fileStorePather fileStorePather) fileStoreQueries {
	return fileStoreQueries{fileStorePather: fileStorePather}
}

var ErrFileNotFound = errors.New("file not found")

func (fileStoreQueries fileStoreQueries) download(filename string) (string, error) {
	filepath := fileStoreQueries.fileStorePather.getFilePath(filename)

	existsResult, err := filesystem.Exists(filepath)
	if err != nil {
		return "", tracerr.Errorf("check exists error: %w", err)
	}
	if existsResult == filesystem.ExistsResultNotFound || existsResult == filesystem.ExistsResultDir {
		return "", tracerr.Wrap(ErrFileNotFound)
	}

	return filepath, nil
}
