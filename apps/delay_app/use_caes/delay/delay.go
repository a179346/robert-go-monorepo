package delay_use_case

import (
	"github.com/a179346/robert-go-monorepo/packages/roberthttp"
)

type DelayUseCase struct{}

func New() DelayUseCase {
	return DelayUseCase{}
}

func (u DelayUseCase) NewGroup() roberthttp.RouterGroup {
	group := roberthttp.NewGroup("/delay")

	group.Handle("GET /{ms}", delayHandler)

	return group
}
