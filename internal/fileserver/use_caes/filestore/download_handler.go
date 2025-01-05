package filestore_use_case

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

var ErrFileNotFound = errors.New("file not found")

func (fs FileStoreUseCase) downloadHandler(c *roberthttp.Context) {
	filename := c.Req.URL().Query().Get("filename")
	filepath, err := downloadQuery(fs.fileStorePather, filename)
	if err != nil {
		var err2 error
		if errors.Is(err, ErrFileNotFound) {
			err2 = c.Res.WriteError(http.StatusBadRequest, err.Error(), nil)
		} else {
			err2 = c.Res.WriteError(http.StatusInternalServerError, "Something went wrong", nil)
		}
		if err2 != nil {
			log.Printf("Error writing response: %v", err2)
		}
		return
	}

	c.Res.SetHeader("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	c.Res.ServeFile(c.Req, filepath)
}

func downloadQuery(fileStorePather fileStorePather, filename string) (string, error) {
	filepath := fileStorePather.getFilePath(filename)

	existsResult, err := filesystem.Exists(filepath)
	if err != nil {
		return "", err
	}
	if existsResult == filesystem.ExistsResultNotFound || existsResult == filesystem.ExistsResultDir {
		return "", ErrFileNotFound
	}

	return filepath, nil
}
