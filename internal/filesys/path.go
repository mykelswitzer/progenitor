package filesys

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

func addTrailingFileSeparator(s string) string {
	if !strings.HasSuffix(s, FilePathSeparator) {
		s = s + FilePathSeparator
	}
	return s
}

func SetBasePath(path string) afero.Fs {

	if path[:1] != "/" {
		base, err := os.Getwd()
		if err != nil {
			log.Println("Relative path provided, unable to determine root.")
			os.Exit(1)
		}
		path = filepath.Join(base, path)
	}

	return afero.NewBasePathFs(afero.NewOsFs(), path)
}

// see Go source code:
// https://github.com/golang/go/blob/f57ebed35132d02e5cf016f324853217fb545e91/src/cmd/go/internal/modload/init.go#L1283
func findModuleRoot(dir string) (roots string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir { // the parent of the root is itself, so we can go no further
			break
		}
		dir = d
	}
	return ""
}

// GetTempDir returns a temporary directory with the given sub path.
func GetTempDir(subPath string, fs afero.Fs) string {
	return afero.GetTempDir(fs, subPath)
}

// DirExists checks if a path exists and is a directory.
func DirExists(path string, fs afero.Fs) (bool, error) {
	return afero.DirExists(fs, path)
}

// IsDir checks if a given path is a directory.
func IsDir(path string, fs afero.Fs) (bool, error) {
	return afero.IsDir(fs, path)
}

// IsEmpty checks if a given path is empty.
func IsEmpty(path string, fs afero.Fs) (bool, error) {
	return afero.IsEmpty(fs, path)
}

// FileContains checks if a file contains a specified string.
func FileContains(filename string, subslice []byte, fs afero.Fs) (bool, error) {
	return afero.FileContainsBytes(fs, filename, subslice)
}

// FileContainsAny checks if a file contains any of the specified strings.
func FileContainsAny(filename string, subslices [][]byte, fs afero.Fs) (bool, error) {
	return afero.FileContainsAnyBytes(fs, filename, subslices)
}

// Exists checks if a file or directory exists.
func Exists(path string, fs afero.Fs) (bool, error) {
	return afero.Exists(fs, path)
}

// AddTrailingSlash adds a trailing Unix styled slash (/) if not already
// there.
func AddTrailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

// GetParentDirFromPath retrieves the parent directory of the final element in a file path.
func GetDirAndParentFromPath(filePath string) (dirName string, parent string) {
	dirName = filepath.Base(filePath)
	parent = ""
	if strings.Contains(filePath, "/") {
		parent = filepath.Base(strings.Replace(filePath, "/"+dirName, "", 1))
	}

	return dirName, parent
}
