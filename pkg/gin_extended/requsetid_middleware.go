package gin_extended

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIdContextKey = "GINEXT-Request-ID"

func RequestIdMiddleware(c *gin.Context) {
	requestId := uuid.New()
	c.Set(requestIdContextKey, requestId)
	c.Next()
}

func GetId(c *gin.Context) (uuid.UUID, bool) {
	return getContextValue[uuid.UUID](c, requestIdContextKey, uuid.Nil)
}
