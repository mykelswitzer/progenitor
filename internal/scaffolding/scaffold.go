package scaffolding

import (
	"log"
	"path/filepath"
	txttmpl "text/template"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	"github.com/spf13/afero"
)

type scaffold interface {
	Init(*config.Config) (*Scaffold, error)
	GetData(*config.Config) templateData
}

type Dir struct {
	Name    string
	SubDirs []Dir
}

func (d *Dir) AddSubDirs(subdirs ...Dir) *Dir {
	d.SubDirs = append(d.SubDirs, subdirs...)
	return d
}

type Scaffold struct {
	Source        scaffold
	Config        *config.Config
	BaseDir       Dir
	Fs            afero.Fs
	TemplatePath  string
	SkipTemplates []string
	ProcessHooks  map[string]func(*Scaffold) error
}

var scaffoldingTypes = map[string]scaffold{
	"go-grpc": goGrpc{},
}

// New will return a scaffold based on the project type
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

// BuildStructure is responsible for creating the project
// folder structure in the local directory
func (s *Scaffold) BuildStructure() error {
	err := createDirs(s.BaseDir.SubDirs, s.Fs)
	if err != nil {
		return err
	}

	if _, ok := s.ProcessHooks["postBuildStructure"]; ok {
		log.Println("Running postBuildStructure")
		err := s.ProcessHooks["postBuildStructure"](s)
		if err != nil {
			return err
		}
	}
	return nil
}

// createDirs recursively reads the project map from the
// scaffold and creates the directories as needed
func createDirs(dirs []Dir, parent afero.Fs) error {

	for _, dir := range dirs {
		if err := parent.MkdirAll(dir.Name, 0777); err != nil {
			return errors.Wrap(err, "Failed to create dir")
		}
		if len(dir.SubDirs) > 0 {
			parentDir := afero.NewBasePathFs(parent, dir.Name)
			err := createDirs(dir.SubDirs, parentDir)
			if err != nil {
				return errors.Wrap(err, "Failed to create dir")
			}
		}
	}

	return nil
}

// BuildFiles sources the templates from the repo, then executes them to
// build the project files in the local directory
func (s *Scaffold) BuildFiles(token string) error {

	base := "github.com/caring/progenitor/internal/templates"
	templates, err := getLatestTemplates(token, filepath.Join(base, s.TemplatePath), s.SkipTemplates, s.Fs)
	if err != nil {
		return err
	}
	if err = s.populateFiles(templates); err != nil {
		return err
	}

	if _, ok := s.ProcessHooks["postBuildFiles"]; ok {
		log.Println("Running postBuildFiles")
		err := s.ProcessHooks["postBuildFiles"](s)
		if err != nil {
			return err
		}
	}
	return nil
}

// populateFiles ranges over the templates and passes in the data
// and executes the template
func (s *Scaffold) populateFiles(templates map[string]*txttmpl.Template) error {

	data := s.Source.GetData(s.Config)

	for path, tmpl := range templates {
		f, err := OpenFileForWriting(s.Fs, trimTmplSuffix(path))
		if err != nil {
			return errors.Wrap(err, "Unable to open file for writing")
		}
		log.Println("Executing template", tmpl)
		err = tmpl.Execute(f, data)
		if err != nil {
			return errors.Wrap(err, "Unable to execute template")
		}
	}

	return nil
}
