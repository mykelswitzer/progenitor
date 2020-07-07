
package cmd


import (
	"errors"
	"io/ioutil"
	"os"
)
import "github.com/manifoldco/promptui"


func promptProjectName() (string, error) {

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

	return prompt.Run()

}

func promptProjectDir() (string, error) {

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

	return prompt.Run()

}

func IsValid(fp string) bool {
  // Check if file already exists
  if _, err := os.Stat(fp); err == nil {
    return true
  }

  // Attempt to create it
  var d []byte
  if err := ioutil.WriteFile(fp, d, 0644); err == nil {
    os.Remove(fp) // And delete it
    return true
  }

  return false
}