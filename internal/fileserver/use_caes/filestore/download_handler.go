package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
	"github.com/a179346/robert-go-monorepo/pkg/gohf/gohf_responses"
)

var ErrFileNotFound = errors.New("file not found")

func (fs FileStoreUseCase) downloadHandler(c *gohf.Context) gohf.Response {
	filename := c.Req.GetQuery("filename")
	filepath, err := downloadQuery(fs.fileStorePather, filename)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return gohf_responses.NewErrorResponse(http.StatusNotFound, err)
		}
		return gohf_responses.NewErrorResponse(http.StatusInternalServerError, errors.New("Something went wrong"))
	}

	c.Res.SetHeader("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	c.Res.ServeFile(c.Req, filepath)
	return gohf_responses.NewDummyResponse()
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
