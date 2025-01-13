package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

var ErrFileNotFound = errors.New("file not found")

func (fs FileStoreUseCase) downloadHandler(c *gohf.Context) gohf.Response {
	filename := c.Req.GetQuery("filename")
	filepath, err := fs.fileStoreQueries.download(filename)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return response.Error(http.StatusNotFound, err)
		}
		return response.Error(http.StatusInternalServerError, errors.New("Something went wrong"))
	}

	c.ResHeader().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	return response.ServeFile(filepath)
}
