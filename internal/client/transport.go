package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ploicloud/cli/internal/auth"
	"github.com/ploicloud/cli/internal/config"
)

type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	refresher  Refresher
}

type Refresher interface {
	Refresh(ctx context.Context, t *auth.Token) (*auth.Token, error)
}

type noopRefresher struct{}

func (noopRefresher) Refresh(_ context.Context, _ *auth.Token) (*auth.Token, error) {
	return nil, fmt.Errorf("no refresher configured")
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{},
		refresher:  noopRefresher{},
	}
}

func (c *Client) WithRefresher(r Refresher) *Client {
	c.refresher = r
	return c
}

type Request struct {
	Method     string
	PathTmpl   string
	PathParams map[string]string
	Query      map[string]string
	Body       any
}

type Response struct {
	Status int
	Body   []byte
	Header http.Header
}

func (r *Response) JSON() (any, error) {
	if len(r.Body) == 0 {
		return nil, nil
	}
	var v any
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return string(r.Body), nil
	}
	return v, nil
}

func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	tok, err := auth.Load()
	if err != nil {
		return nil, err
	}
	if tok.Expired() {
		if refreshed, rerr := c.refresher.Refresh(ctx, tok); rerr == nil {
			tok = refreshed
			_ = auth.Save(tok)
		}
	}

	resp, err := c.send(ctx, req, tok)
	if err != nil {
		return nil, err
	}
	if resp.Status == http.StatusUnauthorized {
		refreshed, rerr := c.refresher.Refresh(ctx, tok)
		if rerr != nil {
			return resp, nil
		}
		if err := auth.Save(refreshed); err != nil {
			return nil, err
		}
		return c.send(ctx, req, refreshed)
	}
	return resp, nil
}

func (c *Client) send(ctx context.Context, req Request, tok *auth.Token) (*Response, error) {
	u, err := c.buildURL(req)
	if err != nil {
		return nil, err
	}
	var body io.Reader
	if req.Body != nil {
		raw, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(raw)
	}
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, u, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")
	tokType := tok.TokenType
	if tokType == "" {
		tokType = "Bearer"
	}
	httpReq.Header.Set("Authorization", tokType+" "+tok.AccessToken)
	httpReq.Header.Set("User-Agent", UserAgent())

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Status: resp.StatusCode, Body: raw, Header: resp.Header}, nil
}

func (c *Client) buildURL(req Request) (string, error) {
	path := req.PathTmpl
	for k, v := range req.PathParams {
		path = strings.ReplaceAll(path, "{"+k+"}", url.PathEscape(v))
	}
	base := strings.TrimRight(c.cfg.APIURL, "/") + "/api/v1"
	u, err := url.Parse(base + path)
	if err != nil {
		return "", err
	}
	if len(req.Query) > 0 {
		q := u.Query()
		for k, v := range req.Query {
			if v != "" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}
