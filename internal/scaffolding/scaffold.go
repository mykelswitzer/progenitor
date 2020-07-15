package scaffolding

import "errors"
import "github.com/caring/progenitor/internal/config"

type scaffold interface {
	Init(*config.Config) (*Scaffold, error)
}

type Dir struct {
	Name string
	SubDirs []Dir
}

type Scaffold struct {
	Config *config.Config
	BaseDir Dir
	TemplatePath string
}

func (s *Scaffold) BuildStructure() error {
	return nil
}

func (s *Scaffold) BuildFiles() error {
	return nil
}


var scaffoldingTypes = map[string]scaffold {
	"go-grpc": goGrpc{},
}

func New(cfg *config.Config) (*scaffold, error) {

	projectType := cfg.GetString("projectType")
	if scfld, ok := scaffoldingTypes["foo"]; ok {
		scfld.Init(cfg)
    return &scfld, nil
	} 

	return nil, errors.New("project scaffold missing: "+projectType)

}