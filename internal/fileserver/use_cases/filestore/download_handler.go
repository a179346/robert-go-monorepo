package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

func (fs FileStoreUseCase) downloadHandler(c *gohf.Context) gohf.Response {
	filename := c.Req.GetQuery("filename")
	filepath, err := fs.fileStoreQueries.download(filename)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return gohf_extended.NewErrorResponse(
				http.StatusNotFound,
				"file not found",
				err,
				false,
			)
		}
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			err,
			true,
		)
	}

	c.ResHeader().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	return response.ServeFile(filepath)
}
