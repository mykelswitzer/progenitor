
package repo

import (
	"net/http"
	"context"
)
import "golang.org/x/oauth2"
import "github.com/google/go-github/v32/github"

func GithubAuth (token string, ctx context.Context) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, ts)
}

func GithubClient(oauth *http.Client) *github.Client {
	return github.NewClient(oauth)
}