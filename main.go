package main

import (
	_ "embed"
	_ "fmt"
	"log"
	"os"
	"github.com/mykelswitzer/progenitor/cmd"
	"github.com/mykelswitzer/progenitor/pkg/config"
)

//go:embed progenitor.yml
var settingsFile string

func main() {

    var exitCode int
    defer func() {
        os.Exit(exitCode)
    }()

	_, err := config.LoadSettings(settingsFile)
	if err != nil {
        log.Print(err)
        exitCode = 1
        return
	}

  	var cfg *config.Config = config.New()

	cmd.Execute(cfg)
}
