package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/mykelswitzer/progenitor/internal/repo"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/scaffold"
)

const BRANCH_DEV = "development"
const BRANCH_MAIN = "main"

func setupRepo(token string, cfg *config.Config) error {

	ctx := context.Background()

	// r here is the remote github repo
	r, err := repo.New(
		ctx,
		token,
		cfg.GetString(config.CFG_PRJ_TEAM),
		cfg.GetString(config.CFG_PRJ_NAME),
		true,
		"Caring, LLC service for "+cfg.GetString(config.CFG_PRJ_NAME),
		true,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create the repo")
	}

	cfg.Set(config.CFG_PRJ_REPO, r)

	// note `lr` is the locally cloned repo, not the same as `repo` returned from
	// github create, which is remote only, largely  because the github library is
	// mostly around github setting, and less about actually working with git...
	lr, err := repo.Clone(ctx, token, cfg.GetString(config.CFG_PRJ_DIR), r)
	if err != nil {
		return err
	}

	err = repo.CreateBranch(token, lr, BRANCH_DEV)
	if err != nil {
		return err
	}

	err = repo.RequireBranchPRApproval(ctx, token, cfg.GetString(config.CFG_PRJ_NAME), BRANCH_MAIN)
	if err != nil {
		return err
	}

	return nil

}

func commitCodeToRepo(token string, cfg *config.Config, s *scaffold.Scaffold) error {

	ctx := context.Background()

	err := repo.AddAll(token, s.BaseDir.Name, s.Fs)
	if err != nil {
		return err
	}

	err = repo.SetDefaultBranch(ctx, token, cfg.GetString(config.CFG_PRJ_NAME), BRANCH_DEV)
	if err != nil {
		return err
	}

	err = repo.RequireBranchPRApproval(ctx, token, cfg.GetString(config.CFG_PRJ_NAME), BRANCH_DEV)
	if err != nil {
		return err
	}

	return nil
}
