package prompt

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/mykelswitzer/progenitor/pkg/config"
	str "github.com/mykelswitzer/progenitor/pkg/strings"
	"github.com/manifoldco/promptui"
)

func UseDB(cfg *config.Config) error {
	return boolPrompt("Do you need a database?", config.CFG_DB_REQ, cfg)
}

func CoreDBObject(cfg *config.Config) error {

	if cfg.GetBool(config.CFG_DB_REQ) == false {
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
		Default: str.ToPascal(cfg.GetString(config.CFG_PRJ_NAME)),
	}

	name, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing DB object name prompt")
	}

	cfg.Set(config.CFG_DB_MDL, name)

	return nil

}
