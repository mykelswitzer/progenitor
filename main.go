package main

import (
	_ "embed"
	_ "fmt"
	"log"
	"os"
	"github.com/mykelswitzer/progenitor/cmd"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/scaffold"

	// these are where the templates live
	gogrpc "github.com/boatsetter/progenitor-tmpl-go-grpc"
)

//go:embed progenitor.yml
var settingsFile string

var scaffolds scaffold.Scaffolds = map[string]scaffold.ScaffoldDS{
	"go-grpc": gogrpc.ScaffoldDS{},
}


func main() {

	cfg, err := config.New(settingsFile)
	if err != nil {
		log.Println(err.Error())
	     os.Exit(0)
	}

	cmd.Execute(cfg, scaffolds)
}
