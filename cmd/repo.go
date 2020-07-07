package cmd


import (
	"context"
	"log"
	"os"
	"os/signal"
)
import "golang.org/x/oauth2"
import "github.com/google/go-github/v32/github"
import "github.com/go-git/go-git/v5"

func createRepo(token string, name string)  *github.Repository {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var private bool = false
	var description string = "Caring, LLC service for " + name;
	var autoInit bool = true
	r := &github.Repository{Name: &name, Private: &private, Description: &description, AutoInit: &autoInit}
	repo, _, err := client.Repositories.Create(ctx, "caring", r)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Printf("Successfully created new repo: %v\n", repo.GetName())
	return repo

}

func cloneRepo(directory string, repo *github.Repository)  {

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

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
