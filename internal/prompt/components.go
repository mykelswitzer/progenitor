package prompt

import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/mykelswitzer/progenitor/v2/pkg/config"
	"github.com/manifoldco/promptui"
)

func boolPrompt(label string, configFld string, config *config.Config) error {

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
