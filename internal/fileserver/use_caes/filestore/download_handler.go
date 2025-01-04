package filestore_use_case

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

func (fs FileStoreUseCase) downloadHandler(c *roberthttp.Context) {
	filename := c.Req.URL().Query().Get("filename")
	filepath := fs.fileStorePather.getFilePath(filename)

	file, err := downloadQuery(filepath)
	if err != nil {
		err = c.Res.WriteError(http.StatusBadRequest, err.Error(), nil)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}
	defer file.Close()

	c.Res.SetHeader("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	c.Res.SetHeader("Content-Type", "application/octet-stream")
	err = pipe(file, c.Res.GetWriter(), 1024)
	if err != nil {
		log.Printf("Writing file failed: %v", err)
	}
}

func downloadQuery(filepath string) (*os.File, error) {
	return os.Open(filepath)
}

func pipe(reader io.Reader, writer io.Writer, chunksize int) error {
	b := make([]byte, chunksize)
	for {
		n, err := reader.Read(b)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = writer.Write(b[:n])
		if err != nil {
			return err
		}
	}
}
