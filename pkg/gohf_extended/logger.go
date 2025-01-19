package gohf_extended

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

type Logger interface {
	Write(v []byte)
}

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func log(w http.ResponseWriter, req *gohf.Request, status int, body interface{}, err error) {
	if logger == nil {
		return
	}

	type tjson map[string]interface{}

	bodyBytes, _ := BodyValue(req.Context())

	content := tjson{
		"unixMilli": req.GetTimestamp().UnixMilli(),
		"time":      req.GetTimestamp().Format(time.RFC3339),
		"req": tjson{
			"uri":    req.RequestURI(),
			"method": req.Method(),
			"header": req.GetHttpRequest().Header,
			"body":   string(bodyBytes),
		},
		"res": tjson{
			"header": w.Header(),
			"status": status,
			"body":   body,
		},
		"error": tracerr.Sprint(err),
	}

	v, err := json.Marshal(content)
	if err != nil {
		return
	}

	logger.Write(v)
}
