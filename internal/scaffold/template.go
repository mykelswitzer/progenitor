package scaffold

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	txttmpl "text/template"

	"github.com/mykelswitzer/progenitor/internal/filesys"
	"github.com/mykelswitzer/progenitor/internal/repo"
	str "github.com/mykelswitzer/progenitor/pkg/strings"
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	_ "github.com/spf13/afero"
)

const TMPLSFX string = ".tmpl"

type tmplParser struct {
	*txttmpl.Template
	txttmpl.FuncMap
}

// TmplParse parses templates from the given filesystem according to the
// given paths. If tmpl is not nil, the templates will be added to it.
// paths must contain at least one path. All paths must exist in the
// given filesystem.
func tmplParse(fs http.FileSystem, funcs txttmpl.FuncMap, tmpl *txttmpl.Template, paths ...string) (*txttmpl.Template, error) {
	t := tmplParser{Template: tmpl, FuncMap: funcs}
	_, err := parseFiles(fs, t.parse, paths...)
	return t.Template, err
}

func (t *tmplParser) parse(name, content string) error {
	var err error
	if t.Template == nil {
		t.Template = txttmpl.New(name)
	} else {
		t.Template = t.New(name)
	}
	if t.FuncMap != nil {
		t.Funcs(t.FuncMap)
	}
	t.Template, err = t.Parse(content)
	return err
}

func parseFiles(fs http.FileSystem, parse func(name string, content string) error, paths ...string) (map[string]*txttmpl.Template, error) {
	if len(paths) == 0 {
		return nil, errors.New("no paths provided")
	}
	buf := bytes.NewBuffer(nil)
	for _, path := range paths {
		f, err := fs.Open(strings.Trim(path, "/"))
		if err != nil {
			return nil, fmt.Errorf("opening template %s %w", path, err)
		}
		name := filepath.Base(path)
		buf.Reset()
		buf.ReadFrom(f)
		err = parse(name, buf.String())
		if err != nil {
			return nil, fmt.Errorf("opening template %s %w", path, err)
		}
	}
	return nil, nil
}

func trimSuffix(path string) string {
	return strings.TrimSuffix(path, TMPLSFX)
}

func stringInSlice(x string, a []string) bool {
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

func getFileSystemHandle(token string, templatePath string) (http.FileSystem, error) {
	ctx := context.Background()
	oauth := repo.OAuthClient(ctx, token)
	return gitfs.New(ctx, templatePath, gitfs.OptClient(oauth))
}

func readFileSystem(remoteFS http.FileSystem, skipTemplates []string) (map[string][]string, []string) {

	var (
		dirs      = map[string][]string{}
		tmplPaths = []string{}
	)

	walker := fsutil.Walk(remoteFS, "")
	for walker.Step() {

		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		filePath := walker.Path()

		switch walker.Stat().IsDir() {
		case true:
			// if it's a directory we need to add it to the dir path
			dirs = mapDirs(dirs, filePath)

		default:
			// if it's a file we need to add it to the templates
			tmplPaths = mapFiles(tmplPaths, skipTemplates, filePath)

		}
	}

	return dirs, tmplPaths
}

func mapDirs(dirs map[string][]string, filePath string) map[string][]string {
	dirName, parent := filesys.GetDirAndParentFromPath(filePath)
	dirs[dirName] = []string{}

	_, ok := dirs[parent]
	if !ok {
		dirs[parent] = []string{dirName}
	} else {
		dirs[parent] = append(dirs[parent], dirName)
	}

	return dirs
}

func populateStructureFromMap(dirMap map[string][]string, rootKey string) (Dir, error) {
	if _, ok := dirMap[rootKey]; !ok {
		return Dir{}, fmt.Errorf("DirMap missing root key: %s", rootKey)
	}

	newDir := Dir{Name: rootKey}
	for _, dir := range dirMap[rootKey] {
		subDir, err := populateStructureFromMap(dirMap, dir)
		if err != nil {
			return Dir{}, fmt.Errorf("Error running populateStructureFromMap with root key: %s %w", rootKey, err)
		}
		newDir.AddSubDirs(subDir)
	}

	return newDir, nil
}

func mapFiles(tmplPaths []string, skipTemplates []string, filePath string) []string {
	if stringInSlice(filePath, skipTemplates) == false {
		log.Println("Fetching template: ", trimSuffix(filePath))
		tmplPaths = append(tmplPaths, filePath)
	}

	return tmplPaths
}

func populateTemplatesFromMap(tmplPaths []string, remoteFS http.FileSystem) (map[string]*txttmpl.Template, error) {

	templates := map[string]*txttmpl.Template{}
	for _, tmplPath := range tmplPaths {
		tmpl, err := tmplParse(remoteFS, templateFunctions(), nil, tmplPath)
		if err != nil {
			return templates, fmt.Errorf("Unable to parse template %s %w", tmplPath, err)
		}
		templates[tmplPath] = tmpl
	}
	return templates, nil

}
