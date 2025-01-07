package roberthttp_response

import (
	"context"
	"errors"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type TextHttpResponse struct {
	Status int
	Text   string
}

func NewTextResponse(statusCode int, text string) TextHttpResponse {
	return TextHttpResponse{
		Status: statusCode,
		Text:   text,
	}
}

func (r TextHttpResponse) Send(res roberthttp.Response, req *roberthttp.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "text/plain")
	res.SetStatus(r.Status)
	res.Write([]byte(r.Text))
}
