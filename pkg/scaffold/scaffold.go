package scaffold

import (
	"log"
	"runtime/debug"
	"strings"
	txttmpl "text/template"

	"github.com/pkg/errors"
	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/prompt"
	"github.com/spf13/afero"
)

type Scaffolds map[string]ScaffoldDS
func (s Scaffolds) Get(k string) (ScaffoldDS, error) {
	scaffold, ok := s[k]
	if !ok {
		return nil, errors.New("map does not contain scaffold: " + k)
	}
	return scaffold, nil
}

// ScaffoldDS is an interface that each template system
// implements to serve as the datasource for populating
// the scaffold and running the commands against
type ScaffoldDS interface {
	Generate(*config.Config, string, afero.Fs) (*Scaffold, error)
	GetData(*config.Config) TemplateData
	GetPrompts() []prompt.PromptFunc
}

type Dir struct {
	Name    string
	SubDirs []Dir
}

func (d *Dir) AddSubDirs(subdirs ...Dir) *Dir {
	d.SubDirs = append(d.SubDirs, subdirs...)
	return d
}

// Scaffold is a common struct that all template systems
// rely on to run the various scaffolding commands
type Scaffold struct {
	Source        ScaffoldDS
	Config        *config.Config
	BaseDir       Dir
	Fs            afero.Fs
	SkipTemplates []string
	ProcessHooks  map[string]func(*Scaffold) error
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

func getScaffoldTemplatePath(projectType string, withVersion bool) string {

	var (
		repoName string = "progenitor-tmpl-" + projectType
		path     string = "github.com/mykelswitzer/" + repoName + "/template"
		version  string
	)

	if withVersion {
		if mf, ok := debug.ReadBuildInfo(); ok {
			for _, m := range mf.Deps {
				if strings.HasSuffix(m.Path, repoName) {
					version = m.Version
					break
				}
			}
		}
		path += "@tags/" + version
	}

	log.Println("reading scaffolding template files from:" + path)

	return path

}

// BuildFiles sources the templates from the repo, then executes them to
// build the project files in the local directory
func (s *Scaffold) BuildFiles() error {

	path := getScaffoldTemplatePath(s.Config.GetString("projectType"), true)
	templates, err := getLatestTemplates(s.Config.GetSettings().GitHub.Token, path, s.SkipTemplates, s.Fs)
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
		f, err := filesys.OpenFileForWriting(s.Fs, trimSuffix(path))
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
