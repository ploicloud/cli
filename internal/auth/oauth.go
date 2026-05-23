package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ClientID     = "019e5393-312b-732f-8292-7a157d677c1e"
	AuthorizeURL = "https://ploi.cloud/oauth/authorize"
	TokenURL     = "https://ploi.cloud/oauth/token"
	RevokeURL    = "https://ploi.cloud/oauth/token"
	Scope        = "*"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Error        string `json:"error"`
	ErrorDesc    string `json:"error_description"`
}

func ExchangeCode(ctx context.Context, code, verifier, redirectURI string) (*Token, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {ClientID},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"code_verifier": {verifier},
	}
	return postToken(ctx, form)
}

func RefreshToken(ctx context.Context, refresh string) (*Token, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {ClientID},
		"refresh_token": {refresh},
		"scope":         {Scope},
	}
	return postToken(ctx, form)
}

func postToken(ctx context.Context, form url.Values) (*Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tr := &tokenResponse{}
	if err := json.Unmarshal(raw, tr); err != nil {
		return nil, fmt.Errorf("decode token response (status %d): %w", resp.StatusCode, err)
	}
	if tr.Error != "" {
		msg := tr.Error
		if tr.ErrorDesc != "" {
			msg = msg + ": " + tr.ErrorDesc
		}
		return nil, fmt.Errorf("oauth error: %s", msg)
	}
	if tr.AccessToken == "" {
		return nil, fmt.Errorf("oauth response missing access_token (status %d): %s", resp.StatusCode, string(raw))
	}
	t := &Token{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		TokenType:    tr.TokenType,
		Scope:        tr.Scope,
	}
	if tr.ExpiresIn > 0 {
		t.ExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	}
	return t, nil
}
