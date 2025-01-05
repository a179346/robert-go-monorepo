package filestore_use_case

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

func (fs FileStoreUseCase) uploadHandler(c *roberthttp.Context) {
	uploadId := c.Req.PathValue("id")

	writerError := func(statusCode int, message string) {
		err := c.Res.WriteError(statusCode, message, nil)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}

	blob, _, err := c.Req.FormFile("blob")
	if err != nil {
		writerError(http.StatusBadRequest, err.Error())
		return
	}
	defer blob.Close()

	offset, err := strconv.Atoi(c.Req.FormValue("offset"))
	if err != nil || offset < 0 {
		writerError(http.StatusBadRequest, "offset should be non-negative integer")
		return
	}

	length, err := strconv.Atoi(c.Req.FormValue("length"))
	if err != nil || length < 0 {
		writerError(http.StatusBadRequest, "length should be non-negative integer")
		return
	}

	isLastChunk := false
	if c.Req.FormValue("isLastChunk") == "true" {
		isLastChunk = true
	} else if c.Req.FormValue("isLastChunk") != "false" {
		writerError(http.StatusBadRequest, "isLastChunk should be true of false")
		return
	}

	filename := c.Req.GetHeader("filename")
	if filename == "" {
		writerError(http.StatusBadRequest, "filename is required")
		return
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
		writerError(http.StatusInternalServerError, err.Error())
		return
	}

	c.Res.SetStatus(http.StatusOK)
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
