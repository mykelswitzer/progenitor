package cmd

import (
	"context"

	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/caring/progenitor/internal/repo"
	"github.com/caring/progenitor/internal/scaffolding"
)

func setupRepo(token string, config *config.Config) error {

	ctx := context.Background()

	// r here is the remote github repo
	r, err := repo.New(
		ctx,
		token,
		config.GetString("projectTeam"),
		config.GetString("projectName"),
		true,
		"Caring, LLC service for "+config.GetString("projectName"),
		true,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create the repo")
	}

	config.Set("projectRepo", r)

	// note `lr` is the locally cloned repo, not the same as `repo` returned from
	// github create, which is remote only, largely  because the github library is
	// mostly around github setting, and less about actually working with git...
	lr, err := repo.Clone(ctx, token, config.GetString("projectDir"), r)
	if err != nil {
		return err
	}

	err = repo.CreateBranch(token, lr, "development")
	if err != nil {
		return err
	}

	err = repo.RequireBranchPRApproval(ctx, token, config.GetString("projectName"), "main")
	if err != nil {
		return err
	}

	return nil

}

func commitCodeToRepo(token string, config *config.Config, s *scaffolding.Scaffold) error {

	ctx := context.Background()

	err := repo.AddAll(token, s.BaseDir.Name, s.Fs)
	if err != nil {
		return err
	}

	err = repo.SetDefaultBranch(ctx, token, config.GetString("projectName"), "development")
	if err != nil {
		return err
	}

	err = repo.RequireBranchPRApproval(ctx, token, config.GetString("projectName"), "development")
	if err != nil {
		return err
	}

	return nil
}
