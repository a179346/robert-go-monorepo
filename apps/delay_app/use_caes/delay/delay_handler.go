package delay_use_case

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a179346/robert-go-monorepo/packages/roberthttp"
)

func (u DelayUseCase) delayHandler(w http.ResponseWriter, req *http.Request) {
	delayMs := req.PathValue("ms")
	d := req.URL.Query().Get("d")

	ms, err := strconv.Atoi(delayMs)
	if err != nil {
		http.Error(w, "Invalid delay", http.StatusBadRequest)
		return
	}

	data, err := delayQuery(req.Context(), ms, d)
	if err != nil {
		log.Printf("Request cancelled: %v", err)
		return
	}

	err = roberthttp.WriteJson(w, http.StatusOK, data)
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
