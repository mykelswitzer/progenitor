package scaffolding

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	txttmpl "text/template"

	"github.com/caring/go-packages/pkg/errors"
	"github.com/caring/progenitor/internal/config"
	rp "github.com/caring/progenitor/internal/repo"
	str "github.com/caring/progenitor/internal/strings"
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

func contains(a []string, x string) bool {
	if len(a) == 0 {
		return false
	}

	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func getLatestTemplates(token string, templatePath string, skipTemplates []string, basePath afero.Fs) (map[string]*txttmpl.Template, error) {

	var templates = map[string]*txttmpl.Template{}

	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	fs, err := gitfs.New(ctx,
		templatePath + "@heads/PRO-72",
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
			//dirpath := filepath.Dir(walker.Path())
			file := walker.Path()
			ex, err := DirExists(filepath.Dir(file), basePath)
			if err != nil {
				return nil, errors.Wrap(err, "Failed reading base path while parsing templates")
			}
			// if the path exists, parse the templates
			if ex && contains(skipTemplates, file) == false {
				log.Println("Fetching template: ", trimTmplSuffix(file))

				tmpl, err := TmplParse(fs, TemplateFunctions(), nil, file)
				if err != nil {
					werr := errors.Wrapf(err, "Unable to parse template %s", file)
					log.Println(werr)
				}
				templates[file] = tmpl
			}
		}
	}

	return templates, nil
}

func TemplateFunctions() txttmpl.FuncMap {
	return txttmpl.FuncMap{
		"tolower":     strings.ToLower,
		"tocamel":     str.ToCamel,
		"topascal":    str.ToPascal,
		"toplural":    str.ToPlural,
		"topackage":   str.ToPackage,
		"tosnakecase": str.ToSnakeCase,
		"toupper":     strings.ToUpper,
	}
}
