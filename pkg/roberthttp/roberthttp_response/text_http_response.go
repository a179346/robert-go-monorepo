package roberthttp_response

import (
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

func (r TextHttpResponse) Send(res roberthttp.Response, _ *roberthttp.Request) {
	res.SetHeader("Content-Type", "text/plain")
	res.SetStatus(r.Status)
	res.GetWriter().Write([]byte(r.Text))
}
