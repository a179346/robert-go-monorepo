package delay_use_case

import (
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
)

type DelayUseCase struct{}

func New() DelayUseCase {
	return DelayUseCase{}
}

func (u DelayUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("GET /{ms}", u.delayHandler)
}
