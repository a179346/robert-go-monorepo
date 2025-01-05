package filestore_use_case

import "path/filepath"

type fileStorePather struct {
	rootPath string
}

func newFileStorePather(rootPath string) fileStorePather {
	return fileStorePather{rootPath: rootPath}
}

func (pather fileStorePather) getFilePath(filename string) string {
	return filepath.Join(pather.rootPath, filename)
}

func (pather fileStorePather) getTempFolder() string {
	return filepath.Join(pather.rootPath, "temp")
}

func (pather fileStorePather) getTempFilePath(id string) string {
	return filepath.Join(pather.getTempFolder(), id)
}
