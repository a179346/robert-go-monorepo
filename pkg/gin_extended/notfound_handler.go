package gin_extended

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		map[string]interface{}{
			"status":  http.StatusNotFound,
			"message": "Not found",
		},
	)
}
