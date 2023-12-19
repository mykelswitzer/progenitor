package progenitor

import (
	"context"

	"github.com/mykelswitzer/progenitor/internal/repo"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/prompt"
	_ "github.com/mykelswitzer/progenitor/pkg/scaffold"
	"github.com/pkg/errors"

	"github.com/spf13/afero"
)

const BRANCH_MAIN = "main"

func setupRepo(cfg *config.Config) error {

	ctx := context.Background()
	token, err := cfg.GetSettings().GitHub.Token(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed getting token to access git filesystem")
	}
	// r here is the remote github repo
	r, err := repo.New(
		ctx,
		cfg.GetSettings().GitHub,
		// cfg.GetString(config.CFG_PRJ_TEAM),
		cfg.GetString(prompt.CfgKeyProjectName),
		true,
		" service for "+cfg.GetString(prompt.CfgKeyProjectName),
		true,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create the repo")
	}

	// cfg.Set(config.CFG_PRJ_REPO, r)

	// note `lr` is the locally cloned repo, not the same as `repo` returned from
	// github create, which is remote only, largely  because the github library is
	// mostly around github setting, and less about actually working with git...
	//lr
	_, err = repo.Clone(ctx, token, cfg.GetString(prompt.CfgKeyProjectDir), r)
	if err != nil {
		return err
	}

	// err = repo.CreateBranch(token, lr, BRANCH_DEV)
	// if err != nil {
	// 	return err
	// }

	// err = repo.RequireBranchPRApproval(ctx, token, cfg.GetString(config.CFG_ProjectName), BRANCH_MAIN)
	// if err != nil {
	// 	return err
	// }

	return nil

}

func commitCodeToRepo(cfg *config.Config, directory string, fileSys afero.Fs) error {

	token, err := cfg.GetSettings().GitHub.Token(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed getting token to access git filesystem")
	}

	err = repo.AddAll(token, directory, fileSys)
	if err != nil {
		return err
	}

	// err = repo.SetDefaultBranch(ctx, token, cfg.GetString(config.CFG_ProjectName), BRANCH_DEV)
	// if err != nil {
	// 	return err
	// }

	// err = repo.RequireBranchPRApproval(ctx, token, cfg.GetString(config.CFG_ProjectName), BRANCH_DEV)
	// if err != nil {
	// 	return err
	// }

	return nil
}
