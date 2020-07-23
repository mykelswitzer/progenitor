package cmd

import (
	_ "io/ioutil"
	_ "os"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/manifoldco/promptui"
)

func promptDb(config *config.Config) error {

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

	return nil

}
