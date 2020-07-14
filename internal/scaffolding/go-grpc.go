package scaffolding

import "log"
import "github.com/spf13/afero"
import "github.com/caring/progenitor/internal/config"

type goGrpc struct {}

func (g goGrpc) Init(cfg *config.Config) (*Scaffold, error) {


	projectDirectory := cfg.Get("projectDir").(afero.Fs)
	log.Print(projectDirectory.Name())

	grpcProject := Scaffold{}

	grpcProject.BaseDir = Dir{Name: projectDirectory.Name()}

	cmdDir := Dir{Name:"cmd"}
	cmdDir.SubDirs = append(cmdDir.SubDirs, Dir{Name:"server"})


	return &grpcProject, nil

}
