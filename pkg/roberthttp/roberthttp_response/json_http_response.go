package roberthttp_response

import (
	"encoding/json"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type JsonHttpResponse[T interface{}] struct {
	Status int
	Data   T
}

func NewJsonResponse[T interface{}](statusCode int, data T) JsonHttpResponse[T] {
	return JsonHttpResponse[T]{
		Status: statusCode,
		Data:   data,
	}
}

func (r JsonHttpResponse[T]) Send(res roberthttp.Response, _ *roberthttp.Request) {
	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(r.Status)
	json.NewEncoder(res.GetWriter()).Encode(r.Data)
}
