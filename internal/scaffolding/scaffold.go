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

func (s *Scaffold) BuildStructure() error {

	

	// 	switch {
	// 	case !isEmpty && !force:
	// 		return errors.New(basepath + " already exists and is not empty. See --force.")

	// 	case !isEmpty && force:
	// 		all := append(dirs, filepath.Join(basepath, "config."+n.configFormat))
	// 		for _, path := range all {
	// 			if exists, _ := Exists(path, fs.Source); exists {
	// 				return errors.New(path + " already exists")
	// 			}
	// 		}
	// 	}
	// }

	// for _, dir := range dirs {
	// 	if err := s.Config.projectDir.MkdirAll(dir, 0777); err != nil {
	// 		return _errors.Wrap(err, "Failed to create dir")
	// 	}
	// }

	return nil
}

func (s *Scaffold) BuildFiles() error {
	return nil
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
