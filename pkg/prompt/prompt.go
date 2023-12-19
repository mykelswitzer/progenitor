package prompt

import (
	"os"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/pkg/errors"
)

type PromptFunc func(cfg *config.Config) error

func BoolPrompt(label string, configFld string, config *config.Config) error {

	output := map[string]bool{"Yes": true, "No": false}
	var keys []string
	for k := range output {
		keys = append(keys, k)
	}

	prompt := promptui.Select{
		Label: label,
		Items: keys,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing '"+label+"' prompt")
	}

	config.Set(configFld, output[result])
	return nil
}

const CfgKeyProjectName = "projectName"

func ProjectName(cfg *config.Config) error {

	validate := func(input string) error {
		if len(input) < 5 {
			return errors.New("Project name must have more than 5 characters")
		}
		re := regexp.MustCompile(`^[a-z\-]+$`)
		if match := re.MatchString(input); !match {
			return errors.New("Project must contain lowercase alphabetical characters with only hyphens as separators.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "What is your project named?",
		Validate: validate,
	}

	name, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing project name prompt")
	}

	cfg.Set(CfgKeyProjectName, strings.ToLower(name))
	return nil
}

const CfgKeyProjectDir = "projectDir"

func ProjectDir(cfg *config.Config) error {

	validate := func(input string) error {
		if input == "" {
			return errors.New("No directory was entered")
		}
		if input[:1] != "/" {
			_, err := os.Getwd()
			if err != nil {
				return errors.New("Relative path provided, unable to determine root.")
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Where should we store this project?",
		Validate: validate,
	}

	dir, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing project directory prompt")
	}

	cfg.Set(CfgKeyProjectDir, dir)
	return nil
}

const CfgKeyUseRemoteRepo = "useRemoteRepo"

func UseRemoteRepository(cfg *config.Config) error {

	prompt := promptui.Select{
		Label: "Would you like to use a remote repository? If yes, a remote repository will be created using the provided credentials in progenitor.yml and the local repository will be pushed to the remote repository.",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing remote repository prompt")
	}

	boolResult := map[string]bool{"Yes": true, "No": false}[result]
	cfg.Set(CfgKeyUseRemoteRepo, boolResult)

	return nil

}
