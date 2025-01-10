package httpclient_extended

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gohf-http/gohf/gohf_responses"
)

func HandleResponse[T any](resp *http.Response, responseObject *T) (*T, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var errResponse gohf_responses.ErrorResponse[interface{}]
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil, err
		}

		return nil, errResponse
	}

	if err := json.Unmarshal(body, &responseObject); err != nil {
		return nil, err
	}

	return responseObject, nil
}
