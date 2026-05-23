package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ploicloud/cli/internal/config"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "ploicloud"
	keyringAccount = "default"
)

var ErrNotAuthenticated = errors.New("not authenticated: run `ploicloud login`")

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	Scope        string    `json:"scope,omitempty"`
}

func (t *Token) Expired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().Add(30 * time.Second).After(t.ExpiresAt)
}

func Save(t *Token) error {
	raw, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := keyring.Set(keyringService, keyringAccount, string(raw)); err == nil {
		return nil
	}
	return saveFile(raw)
}

func Load() (*Token, error) {
	if raw, err := keyring.Get(keyringService, keyringAccount); err == nil {
		t := &Token{}
		if err := json.Unmarshal([]byte(raw), t); err != nil {
			return nil, err
		}
		return t, nil
	}
	raw, err := loadFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotAuthenticated
		}
		return nil, err
	}
	t := &Token{}
	if err := json.Unmarshal(raw, t); err != nil {
		return nil, err
	}
	return t, nil
}

func Clear() error {
	_ = keyring.Delete(keyringService, keyringAccount)
	path, err := fallbackPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func fallbackPath() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

func saveFile(raw []byte) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	path, err := fallbackPath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("write credentials: %w", err)
	}
	return nil
}

func loadFile() ([]byte, error) {
	path, err := fallbackPath()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}
