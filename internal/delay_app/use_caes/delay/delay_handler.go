package delay_use_case

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

func (u DelayUseCase) delayHandler(c *roberthttp.Context) {
	delayMs := c.Req.PathValue("ms")
	d := c.Req.URL().Query().Get("d")

	ms, err := strconv.Atoi(delayMs)
	if err != nil {
		err = c.Res.WriteError(http.StatusBadRequest, "Invalid delay", nil)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	if ms < 0 || ms > 60000 {
		err = c.Res.WriteError(http.StatusBadRequest, "Delay ms should be 0 ~ 60000", nil)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	data, err := delayQuery(c.Req.Context(), ms, d)
	if err != nil {
		log.Printf("Request cancelled: %v", err)
		return
	}

	err = c.Res.WriteJson(http.StatusOK, data)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func delayQuery(ctx context.Context, ms int, d string) (string, error) {
	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
		return d, nil

	case <-ctx.Done():
		return "", ctx.Err()
	}
}
