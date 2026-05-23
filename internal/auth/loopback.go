package auth

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net"
	"net/http"
	"net/url"
	"time"
)

type CallbackResult struct {
	Code  string
	State string
	Error string
}

var RegisteredLoopbackPorts = []int{42173, 42174, 42175}

type Loopback struct {
	listener      net.Listener
	server        *http.Server
	resultCh      chan *CallbackResult
	state         string
	redirectURI   string
}

func StartLoopback(state string) (*Loopback, error) {
	var ln net.Listener
	var port int
	for _, p := range RegisteredLoopbackPorts {
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err == nil {
			ln = l
			port = p
			break
		}
	}
	if ln == nil {
		return nil, fmt.Errorf("all registered loopback ports busy: %v", RegisteredLoopbackPorts)
	}

	lb := &Loopback{
		listener:    ln,
		resultCh:    make(chan *CallbackResult, 1),
		state:       state,
		redirectURI: fmt.Sprintf("http://127.0.0.1:%d/callback", port),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		res := &CallbackResult{
			Code:  q.Get("code"),
			State: q.Get("state"),
			Error: q.Get("error"),
		}
		writeCallbackHTML(w, res)
		select {
		case lb.resultCh <- res:
		default:
		}
	})

	lb.server = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		_ = lb.server.Serve(ln)
	}()
	return lb, nil
}

func (l *Loopback) RedirectURI() string {
	return l.redirectURI
}

func (l *Loopback) Wait(ctx context.Context) (*CallbackResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-l.resultCh:
		if res.Error != "" {
			return res, fmt.Errorf("oauth provider returned error: %s", res.Error)
		}
		if res.State != l.state {
			return res, errors.New("oauth state mismatch")
		}
		if res.Code == "" {
			return res, errors.New("oauth callback missing code")
		}
		return res, nil
	}
}

func (l *Loopback) Close() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = l.server.Shutdown(shutdownCtx)
}

func writeCallbackHTML(w http.ResponseWriter, res *CallbackResult) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	status := "You are signed in. You can close this window."
	if res.Error != "" {
		status = "Sign-in failed: " + html.EscapeString(res.Error)
	} else if res.Code == "" {
		status = "Sign-in failed: missing code."
	}
	fmt.Fprintf(w, `<!doctype html>
<html><head><meta charset="utf-8"><title>Ploi Cloud CLI</title>
<style>body{font-family:-apple-system,BlinkMacSystemFont,sans-serif;background:#0b1020;color:#e6e9f5;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.card{background:#151a32;padding:32px 40px;border-radius:12px;box-shadow:0 8px 30px rgba(0,0,0,.4);text-align:center}h1{margin:0 0 8px;font-size:18px}p{margin:0;opacity:.8;font-size:14px}</style>
</head><body><div class="card"><h1>Ploi Cloud CLI</h1><p>%s</p></div></body></html>`, status)
}

func BuildAuthorizeURL(authURL, clientID, redirectURI, scope string, p *PKCE) (string, error) {
	u, err := url.Parse(authURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", scope)
	q.Set("state", p.State)
	q.Set("code_challenge", p.Challenge)
	q.Set("code_challenge_method", "S256")
	u.RawQuery = q.Encode()
	return u.String(), nil
}
