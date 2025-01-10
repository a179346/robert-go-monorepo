package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gohf-http/gohf"
	"github.com/gohf-http/gohf/gohf_responses"
)

var ErrFileNotFound = errors.New("file not found")

func (fs FileStoreUseCase) downloadHandler(c *gohf.Context) gohf.Response {
	filename := c.Req.GetQuery("filename")
	filepath, err := fs.fileStoreQueries.download(filename)
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
