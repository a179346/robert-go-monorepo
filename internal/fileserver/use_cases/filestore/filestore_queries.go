package filestore_use_case

import "github.com/a179346/robert-go-monorepo/pkg/filesystem"

type fileStoreQueries struct {
	fileStorePather fileStorePather
}

func newFileStoreQueries(fileStorePather fileStorePather) fileStoreQueries {
	return fileStoreQueries{fileStorePather: fileStorePather}
}

func (fileStoreQueries fileStoreQueries) download(filename string) (string, error) {
	filepath := fileStoreQueries.fileStorePather.getFilePath(filename)

	existsResult, err := filesystem.Exists(filepath)
	if err != nil {
		return "", err
	}
	if existsResult == filesystem.ExistsResultNotFound || existsResult == filesystem.ExistsResultDir {
		return "", ErrFileNotFound
	}

	return filepath, nil
}
