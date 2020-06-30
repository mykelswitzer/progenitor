package cmd


import (
	"errors"
	_"fmt"
	_"strconv"
	"github.com/manifoldco/promptui"
)

func promptReponame() (string, error) {

	validate := func(input string) error {
		if len(input) < 5 {
			return errors.New("Service name must have more than 5 characters")
		}
		return nil
	}

  
	prompt := promptui.Prompt{
		Label:    "What is your service named?",
		Validate: validate,
	}


	return prompt.Run()

}