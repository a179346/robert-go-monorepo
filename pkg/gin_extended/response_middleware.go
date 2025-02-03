package gin_extended

import (
	"time"

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

type ApiLogger interface {
	Dispatch(logData ApiLogData)
}

type ApiLogData struct {
	ID          string             `json:"id"`
	App         string             `json:"app"`
	AppVersion  string             `json:"appVersion"`
	Timestamp   int64              `json:"@timestamp"`
	StartUnixMs int64              `json:"startUnixMs"`
	StartTime   string             `json:"startTime"`
	EndUnixMs   int64              `json:"endUnixMs"`
	EndTime     string             `json:"endTime"`
	ElapsedMs   int64              `json:"elapsedMs"`
	Error       string             `json:"error"`
	Unexpected  bool               `json:"unexpected"`
	Req         ApiLogDataRequest  `json:"req"`
	Res         ApiLogDataResponse `json:"res"`
}

type ApiLogDataRequest struct {
	Uri        string              `json:"uri"`
	Method     string              `json:"method"`
	RemoteAddr string              `json:"remoteAddr"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
}

type ApiLogDataResponse struct {
	Header map[string][]string `json:"header"`
	Status int                 `json:"status"`
	Body   string              `json:"body"`
}

type ApiLoggable interface {
	PrepareApiLog(c *gin.Context) (status int, bodyBytes []byte, logErr error, unexpected bool)
}

func ApiLogMiddleware(appId string, appVersion string, apiLogger ApiLogger) gin.HandlerFunc {
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

		logData := ApiLogData{
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
			Req: ApiLogDataRequest{
				Uri:        c.Request.RequestURI,
				Method:     c.Request.Method,
				RemoteAddr: c.Request.RemoteAddr,
				Header:     c.Request.Header,
				Body:       string(reqBodyBytes),
			},
			Res: ApiLogDataResponse{
				Header: c.Writer.Header(),
				Status: status,
				Body:   string(resBodyBytes),
			},
		}

		apiLogger.Dispatch(logData)
	}
}
