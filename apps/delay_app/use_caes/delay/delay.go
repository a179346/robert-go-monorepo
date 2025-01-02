package delay_use_case

import (
	"net/http"
)

type DelayUseCase struct{}

func New() DelayUseCase {
	return DelayUseCase{}
}

func (u DelayUseCase) AddRoutesTo(mux *http.ServeMux) {
	mux.HandleFunc("GET /delay/{ms}", u.delayHandler)
}
