package scaffold

import (
	"fmt"
	"log"
	txttmpl "text/template"

	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/prompt"
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

func (s *Scaffold) Populate(templateRepoPath *string) error {

	if templateRepoPath == nil {
		orgName := s.Config.GetSettings().GitHub.Organization
		projName := s.Config.GetString("projectType")
		tmplFP := getScaffoldTemplatePath(orgName, projName, true)
		templateRepoPath = &tmplFP
	}

	fsh, err := getFileSystemHandle(s.Config.GetSettings().GitHub.Token, *templateRepoPath)
	if err != nil {
		return fmt.Errorf("Failed initializing git filesystem: %w", err)
	}

	_, tmplPaths := readFileSystem(fsh, s.SkipTemplates, s.Fs) //dirs

	if err = s.buildStructure(); err != nil {
		return err
	}

	templates := map[string]*txttmpl.Template{}
	for _, tmplPath := range tmplPaths {
		tmpl, err := filesys.TmplParse(fsh, templateFunctions(), nil, tmplPath)
		if err != nil {
			werr := fmt.Errorf("Unable to parse template %s %w", tmplPath, err)
			log.Println(werr)
		}
		templates[tmplPath] = tmpl
	}

	if err = s.buildFiles(templates); err != nil {
		return err
	}

	return nil
}

// buildStructure is responsible for creating the project
// folder structure in the local directory
func (s *Scaffold) buildStructure() error {

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

// buildFiles sources the templates from the repo, then executes them to
// build the project files in the local directory
func (s *Scaffold) buildFiles(templates map[string]*txttmpl.Template) error {

	if err := s.populateFiles(templates); err != nil {
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
			return fmt.Errorf("Unable to open file for writing %w", err)
		}
		log.Println("Executing template", tmpl)
		// see https://pkg.go.dev/text/template#Template.Execute
		err = tmpl.Execute(f, data)
		if err != nil {
			return fmt.Errorf("Unable to execute template %w", err)
		}
	}

	return nil
}
