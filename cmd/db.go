package cmd


import (
	_"errors"
	_"io/ioutil"
	_"os"
)
import "github.com/manifoldco/promptui"

func promptDb() (bool, error) {

	output := map[string]bool {"Yes":true,"No":false}
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
		return false, err
	}

	return output[result], err

}
