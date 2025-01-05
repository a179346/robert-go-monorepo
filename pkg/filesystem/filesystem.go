package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type ExistsResult int

const (
	ExistsResultNotFound ExistsResult = 0
	ExistsResultFile     ExistsResult = 1
	ExistsResultDir      ExistsResult = 2
)

func Exists(path string) (ExistsResult, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return ExistsResultDir, nil
		}
		return ExistsResultFile, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return ExistsResultNotFound, nil
	}
	return ExistsResultNotFound, err
}

func EnsureDir(path string) error {
	existsResult, err := Exists(path)
	if err != nil {
		return err
	}

	switch existsResult {

	case ExistsResultDir:
		return nil

	case ExistsResultFile:
		return errors.New("ensure dir failed: path is a file")

	case ExistsResultNotFound:
		return os.MkdirAll(path, os.ModePerm)

	default:
		return fmt.Errorf("unknown exists result: %d", existsResult)
	}
}

func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

func OpenWriteOnlyFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY, 0)
}
