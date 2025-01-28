package filestore_use_case

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/a179346/robert-go-monorepo/pkg/filesystem"
	"github.com/ztrue/tracerr"
)

type fileStoreCommands struct {
	fileStorePather fileStorePather
}

func newFileStoreCommands(fileStorePather fileStorePather) fileStoreCommands {
	return fileStoreCommands{fileStorePather: fileStorePather}
}

func (fileStoreCommands fileStoreCommands) upload(
	uploadId string,
	blob multipart.File,
	offset int,
	length int,
	isLastChunk bool,
	filename string,
) error {
	tempFolderPath := fileStoreCommands.fileStorePather.getTempFolder()
	err := filesystem.EnsureDir(tempFolderPath)
	if err != nil {
		return tracerr.Errorf("ensure dir error: %w", err)
	}

	tempFilePath := fileStoreCommands.fileStorePather.getTempFilePath(uploadId)

	chunkData := chunk{
		blob:   blob,
		offset: offset,
		length: length,
	}
	err = chunkData.write(tempFilePath)
	if err != nil {
		return tracerr.Errorf("chunk write error: %w", err)
	}
	if !isLastChunk {
		return nil
	}

	dstFilepath := fileStoreCommands.fileStorePather.getFilePath(filename)
	err = filesystem.MoveFile(tempFilePath, dstFilepath)
	if err != nil {
		return tracerr.Errorf("move file error: %w", err)
	}
	return nil
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
