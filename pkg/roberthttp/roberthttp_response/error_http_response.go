package roberthttp_response

import (
	"encoding/json"
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

func (e ErrorHttpResponse[T]) Error() string {
	return fmt.Sprintf("default http error [%d]: %s", e.Status, e.Message)
}

func (e ErrorHttpResponse[T]) Send(res roberthttp.Response, _ *roberthttp.Request) {
	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(e.Status)
	json.NewEncoder(res.GetWriter()).Encode(e)
}
