package repo

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/mykelswitzer/progenitor/pkg/config"

	"github.com/go-git/go-git/v5"
	_ "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func New(ctx context.Context, ghs config.GitHubSettings, name string, private bool, description string, autoInit bool) (*github.Repository, error) {

	client, err := ghs.Client(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a client to access git")
	}

	mainBranch := "main"

	r := &github.Repository{
		Name:         &name,
		Private:      &private,
		Description:  &description,
		AutoInit:     &autoInit,
		MasterBranch: &mainBranch}
	repo, _, err := client.Repositories.Create(ctx, ghs.Organization, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the repo")
	}

	// opts := &github.TeamAddTeamRepoOptions{Permission: "maintain"}
	// resp, err := client.Teams.AddTeamRepoBySlug(ctx, "mykelswitzer", "engineers", "mykelswitzer", *repo.Name, opts)
	// if err != nil {
	// 	log.Println(err, resp)
	// }

	// opts = &github.TeamAddTeamRepoOptions{Permission: "admin"}
	// resp, err = client.Teams.AddTeamRepoBySlug(ctx, "mykelswitzer", team, "mykelswitzer", *repo.Name, opts)
	// if err != nil {
	// 	log.Println(err, resp)
	// }

	return repo, err

}

func Clone(ctx context.Context, token string, directory string, repo *github.Repository) (*git.Repository, error) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel when we are finished

	go func() {
		<-stop
		log.Printf("\nSignal detected, canceling operation...")
		cancel()
	}()

	cloned, err := git.PlainCloneContext(ctx, directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "mykelswitzer-engineering",
			Password: token,
		},
		URL:      *repo.CloneURL,
		Progress: os.Stdout,
	})

	// note that if err is nil, WithStack also returns nil
	return cloned, errors.Wrap(err, "failed to clone repository")

}

func CreateBranch(token string, repo *git.Repository, name string) error {

	// get the repo worktree
	w, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to access repository worktree")
	}

	opts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
	}

	err = w.Checkout(opts)
	return errors.Wrap(err, "failed to checkout new branch")

}

// func RequireBranchPRApproval(ctx context.Context, token string, repoName string, branchName string) error {
// 	oauth := GithubAuth(token, ctx)
// 	client := GithubClient(oauth)

// 	allowDeletions := false
// 	preq := &github.ProtectionRequest{
// 		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
// 			RequireCodeOwnerReviews:      true,
// 			RequiredApprovingReviewCount: 2,
// 		},
// 		AllowDeletions: &allowDeletions,
// 	}

// 	_, _, err := client.Repositories.UpdateBranchProtection(ctx, "mykelswitzer", repoName, branchName, preq)

// 	return errors.Wrap(err, "failed to setup pr approval requirements for branch "+branchName)
// }

// func SetDefaultBranch(ctx context.Context, token string, repoName string, branchName string) error {

// 	oauth := GithubAuth(token, ctx)
// 	client := GithubClient(oauth)
// 	_, _, err := client.Repositories.Edit(ctx, "mykelswitzer", repoName, &github.Repository{DefaultBranch: &branchName})
// 	return errors.Wrap(err, "failed to set default branch as  "+branchName)
// }

func AddAll(token string, directory string, fs afero.Fs) error {

	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	if err != nil {
		return errors.Wrap(err, "failed to open repository")
	}

	// get the repo worktree
	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to access repository worktree")
	}

	// add all files in the project
	walkAdder := func(path string, f os.FileInfo, err error) error {
		// except the .git dir, which would break things haha
		if f.IsDir() && f.Name() == ".git" {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		_, err = w.Add(path)
		return nil
	}
	afero.Walk(fs, "", walkAdder)

	// check the status
	status, err := w.Status()
	if err != nil {
		return errors.Wrap(err, "failed git status")
	}
	log.Println(status)

	// build the commit
	commit, err := w.Commit("progenitor automated initial commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name: "Progenitor",
			When: time.Now(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed worktree commit")
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		return errors.Wrap(err, "failed commit")
	}
	log.Println(obj)

	// push
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "mykelswitzer-engineering",
			Password: token,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed push")
	}

	return nil

}
