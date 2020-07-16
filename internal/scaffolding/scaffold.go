package scaffolding

import (
	"net/http"
	"errors"
	"log"
	"context"
	"path/filepath"
)
import (
	"github.com/caring/progenitor/internal/config"
	rp "github.com/caring/progenitor/internal/repo"
)
import "github.com/spf13/afero"
import "github.com/posener/gitfs"

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

func (s *Scaffold) BuildFiles(token string) error {
	// get our oauth client
	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	// set up the base template path
	base := "github.com/caring/progenitor/internal/templates"
	templatePath := filepath.Join(base, s.TemplatePath)
	// populate files recursively fetches templates from the 
	// directory, then populates them locally
	populateFiles(ctx, s.BaseDir.SubDirs, s.Fs, oauth, templatePath)

	return nil
}

func populateFiles(ctx context.Context, dirs []Dir, filePath afero.Fs, oauth *http.Client, templatePath string) {
	for _, dir := range dirs {
		getLatestTemplates(ctx, oauth, templatePath)
		// then write te templates to file...
		if len(dir.SubDirs) > 0 {
			filePath := afero.NewBasePathFs(filePath, dir.Name)
			templatePath := filepath.Join(templatePath, dir.Name)
			populateFiles(ctx, dir.SubDirs, filePath, oauth, templatePath)
		}
	}
	//return nil
}

func getLatestTemplates(ctx context.Context, oauth *http.Client, templatePath string) {

	// pull down the latest templates
	fs, err := gitfs.New(ctx,
		templatePath,
		gitfs.OptClient(oauth))
	if err != nil {
		log.Fatalf("Failed initializing git filesystem: %s.", err)
	}

	log.Println(fs)

}
