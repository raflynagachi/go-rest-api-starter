package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Get(ctx context.Context, path string, dst interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return errors.Wrap(err, "client.Get")
	}
	return c.do(req, dst)
}

func (c *Client) Post(ctx context.Context, path string, body, dst interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "client.Post.Marshal")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrap(err, "client.Post.NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, dst)
}

func (c *Client) do(req *http.Request, dst interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.do")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.Errorf("upstream returned status %d", resp.StatusCode)
	}
	if dst == nil {
		return nil
	}
	return errors.Wrap(json.NewDecoder(resp.Body).Decode(dst), "client.do.Decode")
}
