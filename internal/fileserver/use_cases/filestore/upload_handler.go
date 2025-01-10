package filestore_use_case

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gohf-http/gohf"
	"github.com/gohf-http/gohf/gohf_responses"
)

func (fs FileStoreUseCase) uploadHandler(c *gohf.Context) gohf.Response {
	uploadId := c.Req.PathValue("id")

	blob, _, err := c.Req.FormFile("blob")
	if err != nil {
		return gohf_responses.NewErrorResponse(http.StatusBadRequest, err)
	}
	defer blob.Close()

	offset, err := strconv.Atoi(c.Req.FormValue("offset"))
	if err != nil || offset < 0 {
		return gohf_responses.NewErrorResponse(http.StatusBadRequest, errors.New("offset should be non-negative integer"))
	}

	length, err := strconv.Atoi(c.Req.FormValue("length"))
	if err != nil || length < 0 {
		return gohf_responses.NewErrorResponse(http.StatusBadRequest, errors.New("length should be non-negative integer"))
	}

	isLastChunk := false
	if c.Req.FormValue("isLastChunk") == "true" {
		isLastChunk = true
	} else if c.Req.FormValue("isLastChunk") != "false" {
		return gohf_responses.NewErrorResponse(http.StatusBadRequest, errors.New("isLastChunk should be either true or false"))
	}

	filename := c.Req.GetHeader("filename")
	if filename == "" {
		return gohf_responses.NewErrorResponse(http.StatusBadRequest, errors.New("filename is required"))
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
		return gohf_responses.NewErrorResponse(http.StatusInternalServerError, err)
	}

	return gohf_responses.NewTextResponse(http.StatusOK, "OK")
}
