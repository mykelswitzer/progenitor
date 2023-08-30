package repo

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

// GitHubCredsAuth creates and returns a GitHub client authenticated using a personal access token.
// It takes a context and a token string as input and returns a GitHub client.
func GitHubCredsAuth(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauthClient := oauth2.NewClient(ctx, ts)
	return github.NewClient(oauthClient)
}

// GitHubAppAuth creates and returns a GitHub client authenticated as a GitHub App.
// It takes a context, an appID, an installationID, and a private key as input and returns a GitHub client.
func GitHubAppAuth(ctx context.Context, appID, installationID int64, privateKey []byte) *github.Client {

	// Parse the provided private key for RSA signing.
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
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
		Issuer:    strconv.FormatInt(appID, 10),
	}

	// Create a new JWT token with the generated claims and sign it using the private key.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("could not sign jwt token for app: %s", err)
	}

	// Authenticate as a GitHub App and get an installation token.
	appClient := GitHubCredsAuth(ctx, ss)
	appToken, _, err := appClient.Apps.CreateInstallationToken(
		ctx,
		installationID,
		&github.InstallationTokenOptions{})
	if err != nil {
		log.Fatalf("failed to create installation token: %v\n", err)
	}

	// Authenticate using the obtained installation token and return the GitHub client.
	return GitHubCredsAuth(ctx, appToken)
}


