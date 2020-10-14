package cmd

import (
	"context"
	"log"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	rp "github.com/caring/progenitor/internal/repo"
	"github.com/caring/progenitor/internal/scaffolding"
	"github.com/google/go-github/v32/github"
)

func createRepo(token string, config *config.Config) (*github.Repository, error) {

	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	client := rp.GithubClient(oauth)

	var name string = config.GetString("projectName")
	var private bool = true
	var description string = "Caring, LLC service for " + name
	var autoInit bool = true
	r := &github.Repository{Name: &name, Private: &private, Description: &description, AutoInit: &autoInit}
	repo, _, err := client.Repositories.Create(ctx, "caring", r)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create the repo")
	}

	config.Set("projectRepo", repo)

	opts := &github.TeamAddTeamRepoOptions{Permission: "maintain"}
	resp, err := client.Teams.AddTeamRepoBySlug(ctx, "caring", "Engineers", "caring", *repo.Name, opts)
	if err != nil {
		log.Println(err, resp)
	}
	log.Println(resp)

	if err = rp.Clone(token, config.GetString("projectDir"), repo); err != nil {
		return nil, errors.Wrap(err, "Failed to clone the repo")
	}

	return repo, nil

}

func commitCodeToRepo(token string, s *scaffolding.Scaffold) error {
	return rp.AddAll(token, s.BaseDir.Name, s.Fs)
}
