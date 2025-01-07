package roberthttp_response

import (
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type DummyResponse struct {
}

func NewDummyResponse() DummyResponse {
	return DummyResponse{}
}

func (r DummyResponse) Send(_ roberthttp.Response, _ *roberthttp.Request) {}
