package delay_app_client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/httpclient"
	"github.com/a179346/robert-go-monorepo/pkg/httpclient_extended"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp_extended"
)

type Client struct {
	client *httpclient.HTTPClient
}

func New(baseUrl string, client http.Client) *Client {
	return &Client{httpclient.New(baseUrl, client)}
}

type DelayResponse roberthttp_extended.JsonResponse[string]

func (c *Client) Delay(ctx context.Context, ms int, data string) (*DelayResponse, error) {
	resp, err := c.client.Request(ctx, httpclient.RequestOptions{
		Method: http.MethodGet,
		Url:    fmt.Sprintf("/delay/%d", ms),
		Queries: httpclient.Queries{
			"d": data,
		},
	})
	if err != nil {
		return nil, err
	}

	var responseObject DelayResponse
	return httpclient_extended.HandleResponse(resp, &responseObject)
}
