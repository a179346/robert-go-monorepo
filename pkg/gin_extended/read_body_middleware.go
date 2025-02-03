package gin_extended

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ztrue/tracerr"
)

const bodyContextKey = "GINEXT-Body"

func ReadBodyMiddleware(c *gin.Context) {
	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		Error(
			c,
			http.StatusInternalServerError,
			"Something went wrong",
			tracerr.Errorf("read body failed: %w", err),
			true,
		)
		return
	}

	c.Set(bodyContextKey, data)
	c.Next()
}

func GetBody(c *gin.Context) ([]byte, bool) {
	return getContextValue[[]byte](c, bodyContextKey, nil)
}
