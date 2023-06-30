package scaffold

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Dir struct {
	Name    string
	SubDirs []Dir
}

func (d *Dir) AddSubDirs(subdirs ...Dir) *Dir {
	d.SubDirs = append(d.SubDirs, subdirs...)
	return d
}

// createDirs recursively reads the project map from the
// scaffold and creates the directories as needed
func createDirs(dirs []Dir, parent afero.Fs) error {

	for _, dir := range dirs {
		if err := parent.MkdirAll(dir.Name, 0777); err != nil {
			return errors.Wrap(err, "Failed to create dir")
		}
		if len(dir.SubDirs) > 0 {
			parentDir := afero.NewBasePathFs(parent, dir.Name)
			err := createDirs(dir.SubDirs, parentDir)
			if err != nil {
				return errors.Wrap(err, "Failed to create dir")
			}
		}
	}

	return nil
}
