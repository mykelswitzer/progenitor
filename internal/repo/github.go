package repo

import (
	"context"
	"net/http"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func GithubAuth(token string, ctx context.Context) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, ts)
}

func GithubClient(oauth *http.Client) *github.Client {
	return github.NewClient(oauth)
}
