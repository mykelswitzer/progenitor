package scaffolding

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	txttmpl "text/template"
)
import (
	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	rp "github.com/caring/progenitor/internal/repo"
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	"github.com/spf13/afero"
)

type templateData interface {
	Init(config *config.Config) templateData
}

const TMPLSFX string = ".tmpl"

func trimTmplSuffix(path string) string {
	return strings.TrimSuffix(path, TMPLSFX)
}

func getLatestTemplates(token string, templatePath string, basePath afero.Fs) (map[string]*txttmpl.Template, error) {

	var templates = map[string]*txttmpl.Template{}

	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	fs, err := gitfs.New(ctx,
		templatePath,
		gitfs.OptClient(oauth),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed initializing git filesystem")
	}

	walker := fsutil.Walk(fs, "")
	for walker.Step() {

		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if !walker.Stat().IsDir() {
			// get && check path
			dirpath := filepath.Dir(walker.Path())
			log.Println(dirpath)
			ex, err := DirExists(filepath.Dir(walker.Path()), basePath)
			if err != nil {
				return nil, errors.Wrap(err, "Failed reading base path while parsing templates")
			}
			// if the path exists, parse the templates
			if ex {
				log.Println("Fetching template: ", trimTmplSuffix(walker.Path()))

				tmpl, err := TmplParse(fs, TemplateFunctions(), nil, walker.Path())
				if err != nil {
					werr := errors.Wrapf(err, "Unable to parse template %s", walker.Path())
					log.Println(werr)
				}
				templates[walker.Path()] = tmpl
			}
		}
	}

	return templates, nil
}

func TemplateFunctions() txttmpl.FuncMap {
	return txttmpl.FuncMap{
		"tolower":   strings.ToLower,
		"topascal":  ToPascal,
		"topackage": ToPackage,
	}
}

func ToPascal(s string) string {
	a := regexp.MustCompile(`-`)
	words := a.Split(s, -1)
	for index, word := range words {
		words[index] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func ToPackage(s string) string {
	a := regexp.MustCompile(`-`)
	words := a.Split(s, -1)
	return strings.Join(words, "")
}
