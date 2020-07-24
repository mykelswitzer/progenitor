package repo

import (
	"context"
	_ "io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	_ "strings"
	"time"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/afero"
)

func Clone(directory string, repo *github.Repository) error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	go func() {
		<-stop
		log.Printf("\nSignal detected, canceling operation...")
		cancel()
	}()

	_, err := git.PlainCloneContext(ctx, directory, false, &git.CloneOptions{
		URL:      *repo.CloneURL,
		Progress: os.Stdout,
	})

	// note that err is nil, WithStack returns nil
	return errors.Wrap(err, "Failed to clone repository")

}

func AddAll(token string, directory string, fs afero.Fs) error {

	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	if err != nil {
		return errors.Wrap(err, "Failed to open repository")
	}

	// get the repo worktree
	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "Failed to access repository worktree")
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
		return errors.Wrap(err, "Failed git status")
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
		return errors.Wrap(err, "Failed worktree commit")
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		return errors.Wrap(err, "Failed commit")
	}
	log.Println(obj)

	// push
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "caring-engineering",
			Password: token,
		},
	})
	if err != nil {
		return errors.Wrap(err, "Failed push")
	}

	return nil

}
