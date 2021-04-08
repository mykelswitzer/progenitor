package prompt

import (
	"os"
	"regexp"
	"strings"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/manifoldco/promptui"
)

func ProjectTeam(config *config.Config) error {

	prompt := promptui.Select{
		Label: "Which team will own this project?",
		Items: []string{"advisors", "care-recipients", "core-services", "providers"},
	}

	_, team, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing project team prompt")
	}

	config.Set("projectTeam", team)

	return nil

}

func ProjectName(config *config.Config) error {

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

	config.Set("projectName", strings.ToLower(name))

	return nil

}

func ProjectDir(config *config.Config) error {

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

	config.Set("projectDir", dir)

	return nil

}
