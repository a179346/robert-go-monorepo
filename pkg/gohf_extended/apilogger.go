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

var apiLogger ApiLogger

func SetApiLogger(logger ApiLogger) {
	apiLogger = logger
}

var appId string

func SetAppId(id string) {
	appId = id
}

var appVersion string

func SetAppVersion(v string) {
	appVersion = v
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

func log(w http.ResponseWriter, req *gohf.Request, status int, body []byte, err error, unexpected bool) {
	if apiLogger == nil {
		return
	}

	bodyBytes, _ := BodyValue(req.Context())
	requestId, _ := IdValue(req.Context())

	startTime := req.GetTimestamp()
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
		Error:       tracerr.Sprint(err),
		Unexpected:  unexpected,
		Req: ApiLogDataRequest{
			Uri:        req.RequestURI(),
			Method:     req.Method(),
			RemoteAddr: req.RemoteAddr(),
			Header:     req.GetHttpRequest().Header,
			Body:       string(bodyBytes),
		},
		Res: ApiLogDataResponse{
			Header: w.Header(),
			Status: status,
			Body:   string(body),
		},
	}

	apiLogger.Dispatch(logData)
}
