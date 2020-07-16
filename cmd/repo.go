package cmd

import (
	"context"
	"log"
	"os"
)
import (
	"github.com/caring/progenitor/internal/config"
	rp "github.com/caring/progenitor/internal/repo"
)
import "github.com/google/go-github/v32/github"

func createRepo(token string, config *config.Config) *github.Repository {

	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	client := rp.GithubClient(oauth)

	var name string = config.GetString("projectName")
	var private bool = false
	var description string = "Caring, LLC service for " + name
	var autoInit bool = true
	r := &github.Repository{Name: &name, Private: &private, Description: &description, AutoInit: &autoInit}
	repo, _, err := client.Repositories.Create(ctx, "caring", r)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	config.Set("projectRepo", repo)

	rp.Clone(config.GetString("projectDir"), repo)

	return repo

}
