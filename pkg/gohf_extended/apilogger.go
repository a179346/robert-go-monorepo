package gohf_extended

import (
	"net/http"
	"time"

	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

type ApiLogger interface {
	Dispatch(logData ApiLogData)
}

type ApiLogData struct {
	ID          string             `json:"id"`
	App         string             `json:"app"`
	AppVersion  string             `json:"appVersion"`
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
	PrepareApiLog(header http.Header) (status int, bodyBytes []byte, logErr error, unexpected bool)
}

func ApiLogMiddleware(appId string, appVersion string, apiLogger ApiLogger) gohf.HandlerFunc {
	return func(c *gohf.Context) gohf.Response {
		res := c.Next()

		if apiLogger == nil {
			return res
		}

		apiLoggable, ok := res.(ApiLoggable)
		if !ok {
			return res
		}

		status, resBodyBytes, logErr, unexpected := apiLoggable.PrepareApiLog(c.ResHeader())

		reqBodyBytes, _ := BodyValue(c.Req.Context())
		requestId, _ := IdValue(c.Req.Context())

		startTime := c.Req.GetTimestamp()
		endTime := time.Now()
		elapsedMs := endTime.UnixMilli() - startTime.UnixMilli()

		logData := ApiLogData{
			ID:          requestId.String(),
			App:         appId,
			AppVersion:  appVersion,
			StartUnixMs: startTime.UnixMilli(),
			StartTime:   startTime.Format(time.RFC3339),
			EndUnixMs:   endTime.UnixMilli(),
			EndTime:     endTime.Format(time.RFC3339),
			ElapsedMs:   elapsedMs,
			Error:       tracerr.Sprint(logErr),
			Unexpected:  unexpected,
			Req: ApiLogDataRequest{
				Uri:        c.Req.RequestURI(),
				Method:     c.Req.Method(),
				RemoteAddr: c.Req.RemoteAddr(),
				Header:     c.Req.GetHttpRequest().Header,
				Body:       string(reqBodyBytes),
			},
			Res: ApiLogDataResponse{
				Header: c.ResHeader(),
				Status: status,
				Body:   string(resBodyBytes),
			},
		}

		apiLogger.Dispatch(logData)

		return res
	}
}
