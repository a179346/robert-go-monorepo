package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

func (fs FileStoreUseCase) uploadHandler(c *gohf.Context) gohf.Response {
	uploadId := c.Req.PathValue("id")

	blob, _, err := c.Req.FormFile("blob")
	if err != nil {
		return response.Error(http.StatusBadRequest, err)
	}
	defer blob.Close()

	offset, err := strconv.Atoi(c.Req.FormValue("offset"))
	if err != nil || offset < 0 {
		return response.Error(http.StatusBadRequest, errors.New("offset should be non-negative integer"))
	}

	length, err := strconv.Atoi(c.Req.FormValue("length"))
	if err != nil || length < 0 {
		return response.Error(http.StatusBadRequest, errors.New("length should be non-negative integer"))
	}

	isLastChunk := false
	if c.Req.FormValue("isLastChunk") == "true" {
		isLastChunk = true
	} else if c.Req.FormValue("isLastChunk") != "false" {
		return response.Error(http.StatusBadRequest, errors.New("isLastChunk should be either true or false"))
	}

	filename := c.Req.GetHeader("filename")
	if filename == "" {
		return response.Error(http.StatusBadRequest, errors.New("filename is required"))
	}

	err = fs.fileStoreCommands.upload(
		uploadId,
		blob,
		offset,
		length,
		isLastChunk,
		filename,
	)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err)
	}

	return response.Text(http.StatusOK, "OK")
}
