package delay_use_case

import (
	"github.com/gohf-http/gohf/v2"
)

type DelayUseCase struct {
	delayQueries delayQueries
}

func New() DelayUseCase {
	return DelayUseCase{
		delayQueries: newDelayQueries(),
	}
}

func (u DelayUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("GET /{ms}", u.delayHandler)
}
