package filestore_use_case

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp/roberthttp_response"
)

func (fs FileStoreUseCase) uploadHandler(c *roberthttp.Context) roberthttp.HttpResponse {
	uploadId := c.Req.PathValue("id")

	blob, _, err := c.Req.FormFile("blob")
	if err != nil {
		return roberthttp_response.NewErrorResponse(http.StatusBadRequest, err)
	}
	defer blob.Close()

	offset, err := strconv.Atoi(c.Req.FormValue("offset"))
	if err != nil || offset < 0 {
		return roberthttp_response.NewErrorResponse(http.StatusBadRequest, errors.New("offset should be non-negative integer"))
	}

	length, err := strconv.Atoi(c.Req.FormValue("length"))
	if err != nil || length < 0 {
		return roberthttp_response.NewErrorResponse(http.StatusBadRequest, errors.New("length should be non-negative integer"))
	}

	isLastChunk := false
	if c.Req.FormValue("isLastChunk") == "true" {
		isLastChunk = true
	} else if c.Req.FormValue("isLastChunk") != "false" {
		return roberthttp_response.NewErrorResponse(http.StatusBadRequest, errors.New("isLastChunk should be either true or false"))
	}

	filename := c.Req.GetHeader("filename")
	if filename == "" {
		return roberthttp_response.NewErrorResponse(http.StatusBadRequest, errors.New("filename is required"))
	}

	err = uploadCommand(
		fs.fileStorePather,
		uploadId,
		blob,
		offset,
		length,
		isLastChunk,
		filename,
	)
	if err != nil {
		return roberthttp_response.NewErrorResponse(http.StatusInternalServerError, err)
	}

	return roberthttp_response.NewTextResponse(http.StatusOK, "OK")
}

func uploadCommand(
	fileStorePather fileStorePather,
	uploadId string,
	blob multipart.File,
	offset int,
	length int,
	isLastChunk bool,
	filename string,
) error {
	tempFolderPath := fileStorePather.getTempFolder()
	err := filesystem.EnsureDir(tempFolderPath)
	if err != nil {
		return err
	}

	tempFilePath := fileStorePather.getTempFilePath(uploadId)

	chunkData := chunk{
		blob:   blob,
		offset: offset,
		length: length,
	}
	err = chunkData.write(tempFilePath)
	if err != nil {
		return err
	}
	if !isLastChunk {
		return nil
	}

	dstFilepath := fileStorePather.getFilePath(filename)
	return filesystem.MoveFile(tempFilePath, dstFilepath)
}

type chunk struct {
	blob   multipart.File
	offset int
	length int
}

func (chunk chunk) write(filepath string) error {
	file, err := chunk.getFile(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := io.NewOffsetWriter(file, int64(chunk.offset))
	_, err = io.Copy(w, chunk.blob)
	if err != nil {
		return err
	}

	return nil
}

func (chunk chunk) getFile(filepath string) (*os.File, error) {
	if chunk.offset == 0 {
		return filesystem.CreateFile(filepath)
	} else {
		return filesystem.OpenWriteOnlyFile(filepath)
	}
}
