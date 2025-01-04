package httpclient

import (
	"context"
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
	url, err := formatUrl(c.baseUrl, options.Url, options.Queries)
	if err != nil {
		return nil, err
	}

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

func formatUrl(baseUrl, path string, queries Queries) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	if path != "" {
		u = u.JoinPath(path)
	}

	q := u.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
