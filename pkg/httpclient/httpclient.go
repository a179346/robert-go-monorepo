package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTPClient struct {
	baseUrl string
	client  http.Client
}

func New(baseUrl string, client http.Client) *HTTPClient {
	return &HTTPClient{
		baseUrl: baseUrl,
		client:  client,
	}
}

type Queries map[string]string

type Headers map[string]string

type RequestOptions struct {
	Method  string
	Url     string
	Body    io.Reader
	Queries Queries
	Headers Headers
}

func (c *HTTPClient) Request(ctx context.Context, options RequestOptions) (*http.Response, error) {
	url := c.formatUrl(options.Url, options.Queries)

	req, err := http.NewRequest(options.Method, url, options.Body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	if options.Headers != nil {
		for k, v := range options.Headers {
			req.Header.Set(k, v)
		}
	}

	return c.client.Do(req)
}

func (c *HTTPClient) formatUrl(path string, queries Queries) string {
	baseUrl := c.baseUrl
	if len(baseUrl) > 0 && baseUrl[len(baseUrl)-1] == '/' {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	queryString := ""
	for k, v := range queries {
		queryString += fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v))
	}

	url := fmt.Sprintf("%s/%s", baseUrl, path)
	if len(queryString) > 0 {
		url += fmt.Sprintf("?%s", queryString)
	}

	return url
}
