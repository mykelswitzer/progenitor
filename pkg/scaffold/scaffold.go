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
	Source        ScaffoldDS // eventually remove
	Config        *config.Config
	BaseDir       Dir      // eventually remove
	Fs            afero.Fs // eventually remove
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

	// spin up connection to read remote templates
	remoteFS, err := getFileSystemHandle(s.Config.GetSettings().GitHub.Token, *templateRepoPath)
	if err != nil {
		return fmt.Errorf("Failed initializing git filesystem: %w", err)
	}

	// read remote templates into directories and template files
	dirMap, tmplPaths := readFileSystem(remoteFS, s.SkipTemplates)

	// prepare directories for writing
	dirStructure, err := populateStructureFromMap(dirMap, "")
	if err != nil {
		return fmt.Errorf("Failed parsing directory structure %w", err)
	}

	// prepare templates for writing
	templates, err := populateTemplatesFromMap(tmplPaths, remoteFS)
	if err != nil {
		return fmt.Errorf("Failed parsing required templates %w", err)
	}

	// setup local file system root
	localFS := filesys.SetBasePath(s.Config.GetString(prompt.PRJ_DIR))

	// build out directories and call any hooks
	if err = s.buildStructure(dirStructure, localFS); err != nil {
		return err
	}

	// build out files and call any hooks
	if err = s.buildFiles(templates, localFS); err != nil {
		return err
	}

	return nil
}

// buildStructure is responsible for creating the project
// folder structure in the local directory
func (s *Scaffold) buildStructure(scaffoldDir Dir, localFS afero.Fs) error {

	err := createDirs(scaffoldDir.SubDirs, localFS)
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
func (s *Scaffold) buildFiles(templates map[string]*txttmpl.Template, localFS afero.Fs) error {

	if err := s.createFiles(templates, localFS); err != nil {
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
func (s *Scaffold) createFiles(templates map[string]*txttmpl.Template, localFS afero.Fs) error {

	data := s.Config.GetInputs()

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
