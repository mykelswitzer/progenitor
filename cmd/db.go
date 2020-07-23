package cmd

import (
	_ "io/ioutil"
	_ "os"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/manifoldco/promptui"
)

func promptDb() (bool, error) {

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
		return false, errors.Wrap(err, "Error in executing database prompt")
	}

	return output[result], err

}
