package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp/roberthttp_response"
)

var ErrFileNotFound = errors.New("file not found")

func (fs FileStoreUseCase) downloadHandler(c *roberthttp.Context) roberthttp.HttpResponse {
	filename := c.Req.GetQuery("filename")
	filepath, err := downloadQuery(fs.fileStorePather, filename)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return roberthttp_response.NewErrorResponse(http.StatusBadRequest, err)
		}
		return roberthttp_response.NewErrorResponse(http.StatusInternalServerError, errors.New("Something went wrong"))
	}

	c.Res.SetHeader("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	c.Res.ServeFile(c.Req, filepath)
	return roberthttp_response.NewDummyResponse()
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
