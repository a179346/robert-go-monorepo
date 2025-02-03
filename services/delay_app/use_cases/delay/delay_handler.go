package delay_use_case

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a179346/robert-go-monorepo/pkg/gin_extended"
	"github.com/gin-gonic/gin"
	"github.com/ztrue/tracerr"
)

func (u DelayUseCase) delayHandler(c *gin.Context) {
	delayMs := c.Param("ms")
	d := c.Query("d")

	ms, err := strconv.Atoi(delayMs)
	if err != nil {
		gin_extended.Error(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Delay should be integer between 0 and 60000. got: %v", delayMs),
			tracerr.Errorf("parse delay error: %w", err),
			false,
		)
		return
	}

	if ms < 0 || ms > 60000 {
		gin_extended.Error(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Delay should be integer between 0 and 60000. got: %v", ms),
			tracerr.Errorf("Delay should be integer between 0 and 60000. got: %v", ms),
			false,
		)
		return
	}

	data, err := u.delayQueries.getResult(c.Request.Context(), ms, d)
	if err != nil {
		gin_extended.Error(
			c,
			http.StatusInternalServerError,
			"Something went wrong",
			err,
			true,
		)
		return
	}

	gin_extended.JSON(c, http.StatusOK, data)
}
