package scaffold

import (
	"log"
	txttmpl "text/template"

	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/prompt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ScaffoldDS is an interface that each template system
// implements to serve as the datasource for populating
// the scaffold and running the commands against
type ScaffoldDS interface {
	GetName() string
	GetDescription() string
	GetPrompts() []prompt.PromptFunc
	GetData(*config.Config) TemplateData
	Generate(*config.Config, string, afero.Fs) (*Scaffold, error)
}

// Scaffolds is a keyed map of ScaffoldDS interfaces
type Scaffolds []ScaffoldDS

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

// BuildFiles sources the templates from the repo, then executes them to
// build the project files in the local directory
func (s *Scaffold) BuildFiles() error {

	orgName := s.Config.GetSettings().GitHub.Organization
	projName := s.Config.GetString("projectType")
	path := getScaffoldTemplatePath(orgName, projName, true)
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
