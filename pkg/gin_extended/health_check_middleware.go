package gin_extended

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthzHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}
