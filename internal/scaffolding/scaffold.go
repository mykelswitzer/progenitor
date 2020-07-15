package scaffolding

import (
	"errors"
	"log"
)
import "github.com/caring/progenitor/internal/config"
import "github.com/spf13/afero"

type scaffold interface {
	Init(*config.Config) (*Scaffold, error)
}

type Dir struct {
	Name    string
	SubDirs []Dir
}

type Scaffold struct {
	Config       *config.Config
	BaseDir      Dir
	Fs           afero.Fs
	TemplatePath string
}

var scaffoldingTypes = map[string]scaffold{
	"go-grpc": goGrpc{},
}

func New(cfg *config.Config) (*Scaffold, error) {

	projectType := cfg.GetString("projectType")
	if scfld, ok := scaffoldingTypes[projectType]; ok {
		scaffold, err := scfld.Init(cfg)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return scaffold, nil
	}

	return nil, errors.New("project scaffold missing: " + projectType)

}

func (s *Scaffold) BuildStructure() error {
	return createDirs(s.BaseDir.SubDirs, s.Fs)
}

func createDirs(dirs []Dir, parent afero.Fs) error {

	for _, dir := range dirs {
		if err := parent.MkdirAll(dir.Name, 0777); err != nil {
			log.Println("Failed to create dir")
			return err // errors.Wrap(err, "Failed to create dir")
		}
		if len(dir.SubDirs) > 0 {
			parentDir := afero.NewBasePathFs(parent, dir.Name)
			err := createDirs(dir.SubDirs, parentDir)
			if err != nil {
				log.Println("Failed to create dir")
				return err // errors.Wrap(err, "Failed to create dir")
			}
		}
	}
	return nil

}

func (s *Scaffold) BuildFiles() error {
	return nil
}
