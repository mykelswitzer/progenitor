package repo

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

func OAuthClient(ctx context.Context, token string) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, ts)
}
