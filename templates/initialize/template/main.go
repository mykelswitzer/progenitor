package main

import (
	_ "embed"
	"github.com/mykelswitzer/progenitor"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/scaffold"
	"log"
	"os"
	// add imports for your templates below
	//
)

//go:embed progenitor.yml
var settingsFile string

var scaffolds scaffold.Scaffolds = []scaffold.ScaffoldDS{
	// add your initialized template scaffold datasources here,
}

func main() {

	cfg, err := config.New(settingsFile)
	if err != nil {
		log.Println(err.Error())
		os.Exit(0)
	}

	progenitor.Execute(cfg, scaffolds)
}
