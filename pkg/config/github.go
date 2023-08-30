package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var authToken string

type GitHubSettings struct {
	Organization string      `yaml:"organization,omitempty"`
	Creds        GitHubCreds `yaml:"creds,omitempty"`
	App          GitHubApp   `yaml:"app,omitempty"`
}

func (s *GitHubSettings) Client(ctx context.Context) (client *github.Client, err error) {

	if s.UseCreds() {
		client = s.Creds.Auth(ctx)
	}

	if s.UseApp() {
		client, err = s.App.Auth(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to authenitcate to github using the installed application")
		}
	}

	return client, nil
}

func (s *GitHubSettings) Token(ctx context.Context) (string, error) {

	if authToken != "" {
		return authToken, nil
	}

	if s.UseCreds() {
		authToken = s.Creds.PAT
	}

	if s.UseApp() {
		token, err := s.App.GetToken(ctx)
		if err != nil {
			return authToken, fmt.Errorf("failed to create installation token: %w", err)
		}
		authToken = *token.Token
	}

	return authToken, nil
}

func (s *GitHubSettings) IsDefined() bool {
	return s.Organization != "" && (s.Creds.IsDefined() || s.App.IsDefined())
}

func (s *GitHubSettings) UseCreds() bool {
	return s.Creds.IsDefined()
}

func (s *GitHubSettings) UseApp() bool {
	return s.App.IsDefined()
}

type GitHubCreds struct {
	PAT string `yaml:"pat,omitempty"`
}

func (s *GitHubCreds) IsDefined() bool {
	return s.PAT != ""
}

// GitHubCredsAuth creates and returns a GitHub client authenticated using a personal access token.
// It takes a context and a token string as input and returns a GitHub client.
func (c GitHubCreds) Auth(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.PAT})
	return github.NewClient(oauth2.NewClient(ctx, ts))
}

type GitHubApp struct {
	ID           int64  `yaml:"id,omitempty"`
	Key          string `yaml:"key,omitempty"`
	Installation int64  `yaml:"installation,omitempty"`
}

func (s *GitHubApp) IsDefined() bool {
	return s.ID != 0 || s.Key != "" || s.Installation != 0
}

func (a GitHubApp) GetToken(ctx context.Context) (*github.InstallationToken, error) {
	// Parse the provided private key for RSA signing.
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(a.Key))
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %s", err)
	}

	// Generate JWT claims for the GitHub App.
	// GitHub rejects expiry and issue timestamps that are not an integer,
	// while the jwt-go library serializes to fractional timestamps, so we
	// truncate them before passing to jwt-go.
	iss := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iss.Add(10 * time.Minute)
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(iss),
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    strconv.FormatInt(a.ID, 10),
	}

	// Create a new JWT token with the generated claims and sign it using the private key.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("could not sign jwt token for app: %s", err)
	}

	// Authenticate as a GitHub App and get an installation token.
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ss})
	appClient := github.NewClient(oauth2.NewClient(ctx, ts))
	appToken, _, err := appClient.Apps.CreateInstallationToken(
		ctx,
		a.Installation,
		&github.InstallationTokenOptions{})

	return appToken, err
}

// GitHubAppAuth creates and returns a GitHub client authenticated as a GitHub App.
// It takes a context, an appID, an installationID, and a private key as input and returns a GitHub client.
func (a GitHubApp) Auth(ctx context.Context) (*github.Client, error) {

	token, err := a.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation token: %w", err)
	}

	// Authenticate using the obtained installation token and return the GitHub client.
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token.Token})
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}
