package delay_use_case

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v3"
	"github.com/gohf-http/gohf/v3/gohf_responses"
)

func (u DelayUseCase) delayHandler(c *gohf.Context) gohf.Response {
	delayMs := c.Req.PathValue("ms")
	d := c.Req.GetQuery("d")

	ms, err := strconv.Atoi(delayMs)
	if err != nil {
		return gohf_responses.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf("Invalid delay: %s", delayMs),
		)
	}

	if ms < 0 || ms > 60000 {
		return gohf_responses.NewErrorResponse(
			http.StatusBadRequest,
			errors.New("Delay ms should be 0 ~ 60000"),
		)
	}

	data, err := u.delayQueries.getResult(c.Req.Context(), ms, d)
	if err != nil {
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, data)
}
