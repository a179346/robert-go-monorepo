package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type ExistsResult struct{ message string }

type existsResults struct {
	NotFound ExistsResult
	File     ExistsResult
	Dir      ExistsResult
}

var (
	ExistsResults = existsResults{
		NotFound: ExistsResult{"not found"},
		File:     ExistsResult{"file"},
		Dir:      ExistsResult{"dir"},
	}
)

func Exists(path string) (ExistsResult, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return ExistsResults.Dir, nil
		}
		return ExistsResults.File, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return ExistsResults.NotFound, nil
	}
	return ExistsResults.NotFound, err
}

func EnsureDir(path string) error {
	existsResult, err := Exists(path)
	if err != nil {
		return err
	}

	switch existsResult {

	case ExistsResults.Dir:
		return nil

	case ExistsResults.File:
		return errors.New("ensure dir failed: path is a file")

	case ExistsResults.NotFound:
		return os.MkdirAll(path, os.ModePerm)

	default:
		return fmt.Errorf("unknown exists result: %v", existsResult)
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
