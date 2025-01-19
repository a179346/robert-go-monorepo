package delay_use_case

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

func (u DelayUseCase) delayHandler(c *gohf.Context) gohf.Response {
	delayMs := c.Req.PathValue("ms")
	d := c.Req.GetQuery("d")

	ms, err := strconv.Atoi(delayMs)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Sprintf("Delay should be integer between 0 and 60000. got: %v", delayMs),
			tracerr.Errorf("parse delay error: %w", err),
		)
	}

	if ms < 0 || ms > 60000 {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Sprintf("Delay should be integer between 0 and 60000. got: %v", ms),
			tracerr.Errorf("Delay should be integer between 0 and 60000. got: %v", ms),
		)
	}

	data, err := u.delayQueries.getResult(c.Req.Context(), ms, d)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			err,
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, data)
}
