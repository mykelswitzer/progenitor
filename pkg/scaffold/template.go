package scaffold

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	txttmpl "text/template"

	"github.com/mykelswitzer/progenitor/internal/filesys"
	rp "github.com/mykelswitzer/progenitor/internal/repo"
	"github.com/mykelswitzer/progenitor/pkg/config"
	str "github.com/mykelswitzer/progenitor/pkg/strings"
	"github.com/pkg/errors"
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	"github.com/spf13/afero"
)

type TemplateData interface {
	Init(config *config.Config) TemplateData
}

const TMPLSFX string = ".tmpl"

func trimSuffix(path string) string {
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

func templateFunctions() txttmpl.FuncMap {
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

func getScaffoldTemplatePath(orgName string, tmplName string, withVersion bool) string {

	var (
		repoName string = "progenitor-tmpl-" + tmplName
		path     string = fmt.Sprintf("github.com/%s/%s/template", orgName, repoName)
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

	return path
}

func getLatestTemplates(token string, templatePath string, skipTemplates []string, basePath afero.Fs) (map[string]*txttmpl.Template, error) {

	fs, err := getTemplateFileSystem(token, templatePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed initializing git filesystem")
	}

	return readTemplateFileSystem(fs, skipTemplates, basePath)
}

func getTemplateFileSystem(token string, templatePath string) (http.FileSystem, error) {
	ctx := context.Background()
	oauth := rp.GithubAuth(token, ctx)
	return gitfs.New(ctx, templatePath, gitfs.OptClient(oauth))
}

func readTemplateFileSystem(fs http.FileSystem, skipTemplates []string, basePath afero.Fs) (map[string]*txttmpl.Template, error) {

	var (
		dirs      = map[string]Dir{}
		templates = map[string]*txttmpl.Template{}
	)

	walker := fsutil.Walk(fs, "")
	for walker.Step() {

		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		filePath := walker.Path()

		switch walker.Stat().IsDir() {
		case true: // if it's a directory we need to add it to the dir path
			parseDir(dirs, filePath)

		default: // if it's a file we need to add it to the templates
			ex, err := filesys.DirExists(filepath.Dir(filePath), basePath)
			if err != nil {
				return nil, errors.Wrap(err, "Failed reading base path while parsing templates")
			}
			// if the path exists, parse the templates
			if ex && contains(skipTemplates, filePath) == false {
				log.Println("Fetching template: ", trimSuffix(filePath))
				tmpl, err := filesys.TmplParse(fs, templateFunctions(), nil, filePath)
				if err != nil {
					werr := errors.Wrapf(err, "Unable to parse template %s", filePath)
					log.Println(werr)
				}
				templates[filePath] = tmpl
			}
		}
	}

	// we can delete the skip template
	// after the fact delete(map,key)
	return nil, nil
}

func parseDir(dirs map[string]Dir, fPath string) {

	//if _, ok := myMap[]

	//parent := getParentDirFromPath(fPath)

	//hndlDir.AddSubDirs(Dir{Name: fInfo.Name()})
}

func parseTmpl(templates map[string]*txttmpl.Template, fPath string) {

}
