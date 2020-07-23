package scaffolding

import (
	"context"
	"fmt"
	"log"
	"os"
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
)

type templateData interface {
	Init(config *config.Config) templateData
}

const TMPLSFX string = ".tmpl"

func trimTmplSuffix(path string) string {
	return strings.TrimSuffix(path, TMPLSFX)
}

func getLatestTemplates(token string, templatePath string) (map[string]*txttmpl.Template, error) {

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
			log.Println("Attempting to create file: ", trimTmplSuffix(walker.Path()))

			tmpl, err := fsutil.TmplParse(fs, nil, walker.Path())
			if err != nil {
				werr := errors.Wrapf(err, "Unable to parse template %s", walker.Path())
				log.Println(werr)
			}
			templates[walker.Path()] = tmpl
		}
	}

	return templates, nil
}

func getTemplateFunctions() txttmpl.FuncMap {
	return txttmpl.FuncMap{
		"tolower": toLower,
		"tocamel": toCamel,
	}
}

func toLower(s string) string {
	return strings.ToLower(s)
}

func toCamel(s string) string {
	a := regexp.MustCompile(`-`)
	words := a.Split(s, -1)
	for index, word := range words {
		words[index] = strings.Title(word)
	}
	return strings.Join(words, "")
}
