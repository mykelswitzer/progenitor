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

	cfg, err := config.New(settingsFile)
	if err != nil {
		log.Println(err.Error())
	     os.Exit(0)
	}

	cmd.Execute(cfg)
}
