package auth

import "context"

type DefaultRefresher struct{}

func (DefaultRefresher) Refresh(ctx context.Context, t *Token) (*Token, error) {
	if t == nil || t.RefreshToken == "" {
		return nil, ErrNotAuthenticated
	}
	return RefreshToken(ctx, t.RefreshToken)
}
