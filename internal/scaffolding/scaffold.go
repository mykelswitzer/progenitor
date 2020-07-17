package scaffolding

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"fmt"
	"os"
	txttmpl "text/template"
)
import (
	"github.com/caring/progenitor/internal/config"
	rp "github.com/caring/progenitor/internal/repo"
)
import "github.com/spf13/afero"
import (
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
)

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

var localPath string

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
			log.Println("Failed to create dir", err)
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

func (s *Scaffold) BuildFiles(token string) error {
	// get our oauth client
	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	// set up the base template path
	base := "github.com/caring/progenitor/internal/templates"
	templatePath := filepath.Join(base, s.TemplatePath)

	templates, err := getLatestTemplates(ctx, oauth, templatePath)
	if err != nil {
		log.Println(err)
	}

	// populate files recursively fetches templates from the
	// directory, then populates them locally
	populateFiles(s.Fs, s.Config, templates)

	return nil
}

func populateFiles(fs afero.Fs, cfg *config.Config, templates map[string]*txttmpl.Template) {

	for path, tmpl := range templates {

		f, err := OpenFileForWriting(fs, strings.TrimSuffix(path, ".tmpl"))
		if err != nil {
			// handle error
		}

		// Execute the template to the file.
		err = tmpl.Execute(f, cfg)
		if err != nil {
			// handle error
		}

	}

}

func getLatestTemplates(ctx context.Context, oauth *http.Client, templatePath string) (map[string]*txttmpl.Template, error) {

	var templates = map[string]*txttmpl.Template{}

	// pull down the latest templates
	fs, err := gitfs.New(ctx,
		templatePath,
		gitfs.OptClient(oauth),
	)
	if err != nil {
		log.Fatalf("Failed initializing git filesystem: %s.", err)
	}

	walker := fsutil.Walk(fs, "")
	for walker.Step() {

		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if !walker.Stat().IsDir() {

			log.Println("Attempting to create file: ", walker.Path())

			tmpl, err := fsutil.TmplParse(fs, nil, walker.Path())
			if err != nil {
				log.Println("Unable to parse template", err)
				return nil, err // errors.Wrap(err, "Failed to create dir")
			}
			templates[walker.Path()] = tmpl
		}
	}

	return templates, nil
}
