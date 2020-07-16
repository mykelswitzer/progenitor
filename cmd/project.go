package cmd

import (
	"errors"
	"log"
	"os"
)
import "github.com/manifoldco/promptui"
import (
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
		return IsValid(input)
	}

	prompt := promptui.Prompt{
		Label:    "Where should we store this project?",
		Validate: validate,
	}

	dir, err := prompt.Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config.Set("projectDir", dir)

}

func IsValid(fp string) error {

	if fp == "" {
		return errors.New("No directory was entered")
	}

	if fp[:1] != "/" {
		_, err := os.Getwd()
		if err != nil {
			return errors.New("Relative path provided, unable to determine root.")
		}
	}

	return nil

}
