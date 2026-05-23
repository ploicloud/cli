package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
)

func TestNewPKCE_VerifierLength(t *testing.T) {
	p, err := NewPKCE()
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Verifier) < 43 || len(p.Verifier) > 128 {
		t.Errorf("verifier length out of RFC 7636 range: %d", len(p.Verifier))
	}
	if strings.ContainsAny(p.Verifier, "+/=") {
		t.Errorf("verifier must be base64url unpadded: %q", p.Verifier)
	}
	if strings.ContainsAny(p.Challenge, "+/=") {
		t.Errorf("challenge must be base64url unpadded: %q", p.Challenge)
	}
	if p.State == "" {
		t.Error("state must not be empty")
	}
}

func TestNewPKCE_ChallengeMatchesVerifier(t *testing.T) {
	p, err := NewPKCE()
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(p.Verifier))
	expected := base64.RawURLEncoding.EncodeToString(sum[:])
	if expected != p.Challenge {
		t.Errorf("challenge does not match S256(verifier)\n  got:  %s\n  want: %s", p.Challenge, expected)
	}
}

func TestNewPKCE_UniquePerCall(t *testing.T) {
	a, err := NewPKCE()
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewPKCE()
	if err != nil {
		t.Fatal(err)
	}
	if a.Verifier == b.Verifier || a.State == b.State {
		t.Error("PKCE values must be unique per call")
	}
}
