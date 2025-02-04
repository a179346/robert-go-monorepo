package apilog

type Data struct {
	ID          string       `json:"id"`
	App         string       `json:"app"`
	AppVersion  string       `json:"appVersion"`
	Timestamp   int64        `json:"@timestamp"`
	StartUnixMs int64        `json:"startUnixMs"`
	StartTime   string       `json:"startTime"`
	EndUnixMs   int64        `json:"endUnixMs"`
	EndTime     string       `json:"endTime"`
	ElapsedMs   int64        `json:"elapsedMs"`
	Error       string       `json:"error"`
	Unexpected  bool         `json:"unexpected"`
	Req         DataRequest  `json:"req"`
	Res         DataResponse `json:"res"`
}

type DataRequest struct {
	Uri        string              `json:"uri"`
	Method     string              `json:"method"`
	RemoteAddr string              `json:"remoteAddr"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
}

type DataResponse struct {
	Header map[string][]string `json:"header"`
	Status int                 `json:"status"`
	Body   string              `json:"body"`
}

type Logger interface {
	Dispatch(data Data)
}
