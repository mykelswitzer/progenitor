package prompt

import "regexp"
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	str "github.com/caring/progenitor/internal/strings"
	"github.com/manifoldco/promptui"
)

func UseDB(config *config.Config) error {
	return boolPrompt("Do you need a database?", "dbRequired", config)
}

func CoreDBObject(config *config.Config) error {

	if config.GetBool("dbRequired") == false {
		return nil
	}

	validate := func(input string) error {
		if len(input) < 5 {
			return errors.New("Name must have more than 5 characters")
		}
		re := regexp.MustCompile(`^[a-zA-Z]+$`)
		if match := re.MatchString(input); !match {
			return errors.New("DB core object name must contain only lowercase alphabetical characters.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "What is your core DB object named (singular)?",
		Validate: validate,
		// in most cases the core object is named same as service
		Default: str.ToPascal(config.GetString("projectName")),
	}

	name, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing DB object name prompt")
	}

	config.Set("dbModel", name)

	return nil

}
