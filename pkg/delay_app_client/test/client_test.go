package delay_app_client_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/delay_app_client"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
)

func TestClient_Delay(t *testing.T) {
	callDelay := func(ms int, data string) (*delay_app_client.DelayResponse, error) {
		client := delay_app_client.New(
			"http://localhost:8080",
			http.Client{Timeout: 5000 * time.Millisecond},
		)

		ctx := context.Background()

		return client.Delay(ctx, ms, data)
	}

	t.Run("ms: 200 should success", func(t *testing.T) {
		response, err := callDelay(200, "Hello, World!")

		if err != nil {
			t.Error("error should be nil")
			return
		}

		if response == nil || response.Data != "Hello, World!" {
			t.Error("response data should be \"Hello, World!\"")
			return
		}
	})

	t.Run("ms: -1 should return error", func(t *testing.T) {
		res, err := callDelay(-1, "Hello, World!")

		if err == nil {
			t.Error("error should be of type ErrorResponseData")
			return
		}

		e, ok := err.(gohf_extended.ErrorResponseData)
		if !ok {
			t.Error("error should be of type ErrorResponseData")
			return
		}

		expectedErrorMessage := "Delay ms should be 0 ~ 60000"
		if e.Message != expectedErrorMessage {
			t.Errorf("expected error message: \"%s\" received:\"%s\"", expectedErrorMessage, e.Message)
			return
		}

		if res != nil {
			t.Error("response data should be nil")
			return
		}
	})
}
