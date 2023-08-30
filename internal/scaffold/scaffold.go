package scaffold

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	txttmpl "text/template"

	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/spf13/afero"
)

// Scaffold is a common struct that all template systems
// rely on to run the various scaffolding commands
type Scaffold struct {
	Config        *config.Config
	SkipTemplates []string
	ProcessHooks  map[string]func(*Scaffold) error
}

func (s *Scaffold) Init(config *config.Config, skip []string, hooks map[string]func(*Scaffold) error) {
	s.Config = config
	s.SkipTemplates = skip
	s.ProcessHooks = hooks

	version := "(undetermined)"
	if mf, ok := debug.ReadBuildInfo(); ok {
		for _, m := range mf.Deps {
			templateName := "progenitor-tmpl-" + config.GetString("projectType")
			if strings.HasSuffix(m.Path, templateName) {
				version = m.Version
				break
			}
		}
	}
	s.Config.Set("Version", version)
}

func (s *Scaffold) Populate(templateRepoPath *string, localFS afero.Fs) error {

	if templateRepoPath == nil {
		orgName := s.Config.GetSettings().GitHub.Organization
		projName := s.Config.GetString("projectType")
		tmplFP := getScaffoldTemplatePath(orgName, projName, true)
		templateRepoPath = &tmplFP
	}

	token, err := s.Config.GetSettings().GitHub.Token(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting token to access git filesystem: %w", err)
	}

	// spin up connection to read remote templates
	remoteFS, err := getFileSystemHandle(token, *templateRepoPath)
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

type Dir struct {
	Name    string
	SubDirs []Dir
}

func (d *Dir) AddSubDirs(subdirs ...Dir) *Dir {
	d.SubDirs = append(d.SubDirs, subdirs...)
	return d
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

// createDirs recursively reads the project map from the
// scaffold and creates the directories as needed locally
func createDirs(dirs []Dir, parent afero.Fs) error {

	for _, dir := range dirs {
		if err := parent.MkdirAll(dir.Name, 0777); err != nil {
			return fmt.Errorf("Failed to create dir %w", err)
		}
		if len(dir.SubDirs) > 0 {
			parentDir := afero.NewBasePathFs(parent, dir.Name)
			err := createDirs(dir.SubDirs, parentDir)
			if err != nil {
				return fmt.Errorf("Failed to create dir %w", err)
			}
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

	for k, v := range data {
		fmt.Println(k, "=", v)
	}

	for path, tmpl := range templates {
		f, err := filesys.OpenFileForWriting(localFS, trimSuffix(path))
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
