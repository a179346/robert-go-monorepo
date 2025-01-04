package delay_use_case

import (
	"github.com/a179346/robert-go-monorepo/packages/roberthttp"
)

type DelayUseCase struct{}

func New() DelayUseCase {
	return DelayUseCase{}
}

func (u DelayUseCase) HandleGroup(group *roberthttp.RouterGroup) {
	group.Handle("GET /{ms}", delayHandler)
}
