package gohf_responses

import (
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
)

type DummyResponse struct {
}

func NewDummyResponse() DummyResponse {
	return DummyResponse{}
}

func (response DummyResponse) Send(_ gohf.ResponseWriter, _ *gohf.Request) {}
