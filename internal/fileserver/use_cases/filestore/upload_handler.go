package filestore_use_case

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

func (fs FileStoreUseCase) uploadHandler(c *gohf.Context) gohf.Response {
	uploadId := c.Req.PathValue("id")

	blob, _, err := c.Req.FormFile("blob")
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			"failed to read blob",
			tracerr.Errorf("read blob error: %w", err),
		)
	}
	defer blob.Close()

	offset, err := strconv.Atoi(c.Req.FormValue("offset"))
	if err != nil || offset < 0 {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Sprintf("offset should be non-negative integer. got: %v", offset),
			tracerr.Errorf("offset should be non-negative integer. got: %v", offset),
		)
	}

	length, err := strconv.Atoi(c.Req.FormValue("length"))
	if err != nil || length < 0 {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Sprintf("length should be non-negative integer. got: %v", length),
			tracerr.Errorf("length should be non-negative integer. got: %v", length),
		)
	}

	isLastChunk := false
	if c.Req.FormValue("isLastChunk") == "true" {
		isLastChunk = true
	} else if c.Req.FormValue("isLastChunk") != "false" {
		v := c.Req.FormValue("isLastChunk")
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Sprintf("isLastChunk should be either true or false. got: %v", v),
			tracerr.Errorf("isLastChunk should be either true or false. got: %v", v),
		)
	}

	filename := c.Req.GetHeader("filename")
	if filename == "" {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			"filename is required",
			tracerr.New("filename is required"),
		)
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
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			err,
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, "OK")
}
