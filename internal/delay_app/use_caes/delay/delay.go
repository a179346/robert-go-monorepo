package delay_use_case

import (
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type DelayUseCase struct{}

func New() DelayUseCase {
	return DelayUseCase{}
}

func (u DelayUseCase) AppendHandler(router *roberthttp.Router) {
	router.Handle("GET /{ms}", u.delayHandler)
}
