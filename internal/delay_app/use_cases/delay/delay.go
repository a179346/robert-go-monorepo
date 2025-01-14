package delay_use_case

import (
	"github.com/gohf-http/gohf/v6"
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
	router.GET("/{ms}", u.delayHandler)
}
