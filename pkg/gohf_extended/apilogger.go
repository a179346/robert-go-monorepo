package gohf_extended

import (
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

func ApiLogMiddleware(appId string, appVersion string, apiLogger apilog.Logger) gohf.HandlerFunc {
	return func(c *gohf.Context) gohf.Response {
		res := c.Next()

		if apiLogger == nil {
			return res
		}

		apiLoggable, ok := res.(apilog.Loggable)
		if !ok {
			return res
		}

		status, resBodyBytes, logErr, unexpected := apiLoggable.PrepareApiLog(c.ResHeader())

		reqBodyBytes, _ := BodyValue(c.Req.Context())
		requestId, _ := IdValue(c.Req.Context())

		startTime := c.Req.GetTimestamp()
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
				Uri:        c.Req.RequestURI(),
				Method:     c.Req.Method(),
				RemoteAddr: c.Req.RemoteAddr(),
				Header:     c.Req.GetHttpRequest().Header,
				Body:       string(reqBodyBytes),
			},
			Res: apilog.DataResponse{
				Header: c.ResHeader(),
				Status: status,
				Body:   string(resBodyBytes),
			},
		}

		apiLogger.Dispatch(data)

		return res
	}
}
