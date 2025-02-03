package delay_use_case

import "github.com/gin-gonic/gin"

type DelayUseCase struct {
	delayQueries delayQueries
}

func New() DelayUseCase {
	return DelayUseCase{
		delayQueries: newDelayQueries(),
	}
}

func (u DelayUseCase) AppendHandler(router *gin.RouterGroup) {
	router.GET("/:ms", u.delayHandler)
}
