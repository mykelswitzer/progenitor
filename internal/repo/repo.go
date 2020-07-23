package repo

import (
	"context"
	"log"
	"os"
	"os/signal"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v32/github"
)

func Clone(directory string, repo *github.Repository) error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// The context is the mechanism used by go-git, to support deadlines and
	// cancellation signals.
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
	return errors.WithStack(err)

}
