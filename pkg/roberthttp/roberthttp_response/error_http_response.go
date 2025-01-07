package roberthttp_response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type ErrorHttpResponse[T interface{}] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     T      `json:"err"`
}

func NewErrorResponse[TError error](statusCode int, err TError) ErrorHttpResponse[TError] {
	return ErrorHttpResponse[TError]{
		Status:  statusCode,
		Message: err.Error(),
		Err:     err,
	}
}

func (r ErrorHttpResponse[T]) Error() string {
	return fmt.Sprintf("default http error [%d]: %s", r.Status, r.Message)
}

func (r ErrorHttpResponse[T]) Send(res roberthttp.Response, req *roberthttp.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(r.Status)
	json.NewEncoder(res).Encode(r)
}
