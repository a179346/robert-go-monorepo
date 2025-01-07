package roberthttp_response

import (
	"context"
	"errors"
	"io"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type IoHttpResponse struct {
	Status int
	Reader io.Reader
}

func NewIoResponse(statusCode int, reader io.Reader) IoHttpResponse {
	return IoHttpResponse{
		Status: statusCode,
		Reader: reader,
	}
}

func (r IoHttpResponse) Send(res roberthttp.Response, req *roberthttp.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetStatus(r.Status)
	io.Copy(res, r.Reader)
}
