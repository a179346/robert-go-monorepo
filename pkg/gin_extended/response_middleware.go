package gin_extended

import (
	"github.com/gin-gonic/gin"
)

const responseContextKey = "GINEXT-Response"

type Response interface {
	Send(c *gin.Context)
}

func ResponseMiddleware(c *gin.Context) {
	c.Next()

	if response, ok := getResponse(c); ok {
		response.Send(c)
	}
}

func getResponse(c *gin.Context) (Response, bool) {
	return getContextValue[Response](c, responseContextKey, nil)
}

func withResponse(c *gin.Context, res Response) {
	c.Set(responseContextKey, res)
}
