package gohf_extended

import (
	"net/http"
	"time"

	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

type Logger interface {
	Write(logData LogData)
}

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

var appId string

func SetAppId(id string) {
	appId = id
}

type LogData struct {
	ID          string          `json:"id"`
	App         string          `json:"app"`
	StartUnixMs int64           `json:"startUnixMs"`
	StartTime   string          `json:"startTime"`
	EndUnixMs   int64           `json:"endUnixMs"`
	EndTime     string          `json:"endTime"`
	ElapsedMs   int64           `json:"elapsedMs"`
	Error       string          `json:"error"`
	Req         LogDataRequest  `json:"req"`
	Res         LogDataResponse `json:"res"`
}

type LogDataRequest struct {
	Uri    string              `json:"uri"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

type LogDataResponse struct {
	Header map[string][]string `json:"header"`
	Status int                 `json:"status"`
	Body   interface{}         `json:"body"`
}

func log(w http.ResponseWriter, req *gohf.Request, status int, body interface{}, err error) {
	if logger == nil {
		return
	}

	bodyBytes, _ := BodyValue(req.Context())
	requestId, _ := IdValue(req.Context())

	startTime := req.GetTimestamp()
	endTime := time.Now()
	elapsedMs := endTime.UnixMilli() - startTime.UnixMilli()

	logData := LogData{
		ID:          requestId.String(),
		App:         appId,
		StartUnixMs: startTime.UnixMilli(),
		StartTime:   startTime.Format(time.RFC3339),
		EndUnixMs:   endTime.UnixMilli(),
		EndTime:     endTime.Format(time.RFC3339),
		ElapsedMs:   elapsedMs,
		Error:       tracerr.Sprint(err),
		Req: LogDataRequest{
			Uri:    req.RequestURI(),
			Method: req.Method(),
			Header: req.GetHttpRequest().Header,
			Body:   string(bodyBytes),
		},
		Res: LogDataResponse{
			Header: w.Header(),
			Status: status,
			Body:   body,
		},
	}

	logger.Write(logData)
}
