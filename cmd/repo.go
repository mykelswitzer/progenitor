package cmd


import (
	"context"
	"fmt"
	"log"
	"os"
	"errors"
)
import "golang.org/x/oauth2"
import "github.com/manifoldco/promptui"
import "github.com/google/go-github/v32/github"

func promptReponame() (string, error) {

	validate := func(input string) error {
		if len(input) < 5 {
			return errors.New("Project name must have more than 5 characters")
		}
		return nil
	}

  
	prompt := promptui.Prompt{
		Label:    "What is your project named?",
		Validate: validate,
	}


	return prompt.Run()

}

func createRepo(name string) {

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if name == "" {
		log.Fatal("No name: New repos must be given a name")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var private bool = false
	var description string = "Caring, LLC service for " + name;
	r := &github.Repository{Name: &name, Private: &private, Description: &description}
	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())

}