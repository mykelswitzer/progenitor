package cmd


import (
	"errors"
	"os"
  "log"
)
import "github.com/manifoldco/promptui"
import (
	"github.com/caring/progenitor/internal/scaffolding"
	"github.com/caring/progenitor/internal/config"
)


func promptProjectName(config *config.Config) {

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

	name, err := prompt.Run()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	config.Set("projectName", name)

}

func promptProjectDir(config *config.Config) {

	validate := func(input string) error {
		if IsValid(input) == false {
			return errors.New("Directory is invalid or not writeable")
		}
		return nil
	}
  
	prompt := promptui.Prompt{
		Label:    "Where should we store this project?",
		Validate: validate,
	}

	dir, err := prompt.Run()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	config.Set("projectDir", dir)

}

func IsValid(fp string) bool {
  
  scaffolding.SetBasePath(fp)

  return true

}