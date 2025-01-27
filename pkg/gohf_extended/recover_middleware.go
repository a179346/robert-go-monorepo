package gohf_extended

import (
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

func RecoverMiddleware(c *gohf.Context) (res gohf.Response) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = errors.New("unknown panic")
			}

			res = NewErrorResponse(
				http.StatusInternalServerError,
				"Something went wrong",
				err,
				true,
			)
		}
	}()

	return c.Next()
}
