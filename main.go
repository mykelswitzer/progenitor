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

	if _, err := config.LoadSettings(settingsFile); err != nil {
        log.Print(err)
        exitCode = 1
        return
	} 

	cmd.Execute(config.New())
}
