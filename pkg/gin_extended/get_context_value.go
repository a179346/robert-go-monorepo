package gin_extended

import (
	"github.com/gin-gonic/gin"
)

func getContextValue[T interface{}](c *gin.Context, contextKey string, defaultValue T) (T, bool) {
	if v, ok := c.Get(contextKey); ok {
		if value, ok := v.(T); ok {
			return value, true
		}
	}
	return defaultValue, false
}
