package delay_use_case

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp/roberthttp_response"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp_extended"
)

func (u DelayUseCase) delayHandler(c *roberthttp.Context) roberthttp.HttpResponse {
	delayMs := c.Req.PathValue("ms")
	d := c.Req.URL().Query().Get("d")

	ms, err := strconv.Atoi(delayMs)
	if err != nil {
		return roberthttp_response.NewErrorResponse(
			http.StatusBadRequest,
			fmt.Errorf("Invalid delay: %s", delayMs),
		)
	}

	if ms < 0 || ms > 60000 {
		return roberthttp_response.NewErrorResponse(
			http.StatusBadRequest,
			errors.New("Delay ms should be 0 ~ 60000"),
		)
	}

	data, err := delayQuery(c.Req.Context(), ms, d)
	if err != nil {
		log.Printf("Request cancelled: %v", err)
		return nil
	}

	return roberthttp_extended.NewCustomJsonResponse(http.StatusOK, data)
}

func delayQuery(ctx context.Context, ms int, d string) (string, error) {
	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
		return d, nil

	case <-ctx.Done():
		return "", ctx.Err()
	}
}
