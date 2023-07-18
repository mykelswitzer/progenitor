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
	"github.com/posener/gitfs"
	"github.com/posener/gitfs/fsutil"
	_ "github.com/spf13/afero"
)

type TemplateData interface {
	Init(config *config.Config) TemplateData
}

const TMPLSFX string = ".tmpl"

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
	oauth := rp.GithubAuth(token, ctx)
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

// getParentDirFromPath retrieves the parent directory of the final element in a file path.
func getDirAndParentFromPath(filePath string) (dirName string, parent string) {
	dirName = filepath.Base(filePath)
	parent = ""
	if strings.Contains(filePath, "/") {
		parent = filepath.Base(strings.Replace(filePath, "/"+dirName, "", 1))
	}

	return dirName, parent
}

func mapDirs(dirs map[string][]string, filePath string) map[string][]string {
	dirName, parent := getDirAndParentFromPath(filePath)
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
		tmpl, err := filesys.TmplParse(remoteFS, templateFunctions(), nil, tmplPath)
		if err != nil {
			return templates, fmt.Errorf("Unable to parse template %s %w", tmplPath, err)
		}
		templates[tmplPath] = tmpl
	}
	return templates, nil

}