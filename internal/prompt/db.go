package prompt

import "regexp"
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/manifoldco/promptui"
)

func UseDB(config *config.Config) error {

	output := map[string]bool{"Yes": true, "No": false}
	var keys []string
	for k := range output {
		keys = append(keys, k)
	}

	prompt := promptui.Select{
		Label: "Do you need a database?",
		Items: keys,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing database prompt")
	}
	config.Set("requireDb", output[result])

	return CoreDBObject(config)

}

func CoreDBObject(config *config.Config) error {

	validate := func(input string) error {
		if len(input) < 5 {
			return errors.New("Name must have more than 5 characters")
		}
		re := regexp.MustCompile(`^[a-z]+$`)
		if match := re.MatchString(input); !match {
			return errors.New("DB core object name must contain only lowercase alphabetical characters.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "What is your core DB object named (singular)?",
		Validate: validate,
		// in most cases the core object is named same as service
		Default: config.GetString("projectName"),
	}

	name, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "Error in executing DB object name prompt")
	}

	config.Set("dbCoreObject", name)

	return nil

}
