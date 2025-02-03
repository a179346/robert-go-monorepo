package gin_extended

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxBytesMiddleware(size int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, size)

		c.Next()
	}
}
