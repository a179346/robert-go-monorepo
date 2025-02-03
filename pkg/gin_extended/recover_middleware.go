package gin_extended

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware(c *gin.Context) {
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

			Error(
				c,
				http.StatusInternalServerError,
				"Something went wrong",
				err,
				true,
			)
		}
	}()

	c.Next()
}
