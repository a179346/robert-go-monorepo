package gin_extended

import (
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/gin-gonic/gin"
	"github.com/ztrue/tracerr"
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

type ApiLoggable interface {
	PrepareApiLog(c *gin.Context) (status int, bodyBytes []byte, logErr error, unexpected bool)
}

func ApiLogMiddleware(appId string, appVersion string, apiLogger apilog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()
		if apiLogger == nil {
			return
		}
		res, ok := getResponse(c)
		if !ok {
			return
		}

		apiLoggable, ok := res.(ApiLoggable)
		if !ok {
			return
		}

		status, resBodyBytes, logErr, unexpected := apiLoggable.PrepareApiLog(c)

		reqBodyBytes, _ := GetBody(c)
		requestId, _ := GetId(c)

		endTime := time.Now()
		startUnixMs := startTime.UnixMilli()
		endUnixMs := endTime.UnixMilli()
		elapsedMs := endUnixMs - startUnixMs

		data := apilog.Data{
			ID:          requestId.String(),
			App:         appId,
			AppVersion:  appVersion,
			Timestamp:   startUnixMs,
			StartUnixMs: startUnixMs,
			StartTime:   startTime.Format(time.RFC3339),
			EndUnixMs:   endUnixMs,
			EndTime:     endTime.Format(time.RFC3339),
			ElapsedMs:   elapsedMs,
			Error:       tracerr.Sprint(logErr),
			Unexpected:  unexpected,
			Req: apilog.DataRequest{
				Uri:        c.Request.RequestURI,
				Method:     c.Request.Method,
				RemoteAddr: c.Request.RemoteAddr,
				Header:     c.Request.Header,
				Body:       string(reqBodyBytes),
			},
			Res: apilog.DataResponse{
				Header: c.Writer.Header(),
				Status: status,
				Body:   string(resBodyBytes),
			},
		}

		apiLogger.Dispatch(data)
	}
}
