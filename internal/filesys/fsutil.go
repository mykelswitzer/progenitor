package filesys

import (
	"bytes"
	"net/http"
	"path/filepath"
	"strings"
	txttmpl "text/template"

	"github.com/pkg/errors"
)

// TmplParse parses templates from the given filesystem according to the
// given paths. If tmpl is not nil, the templates will be added to it.
// paths must contain at least one path. All paths must exist in the
// given filesystem.
func TmplParse(fs http.FileSystem, funcs txttmpl.FuncMap, tmpl *txttmpl.Template, paths ...string) (*txttmpl.Template, error) {
	t := tmplParser{Template: tmpl, FuncMap: funcs}
	_, err := parseFiles(fs, t.parse, paths...)
	return t.Template, err
}

type tmplParser struct {
	*txttmpl.Template
	txttmpl.FuncMap
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
			return nil, errors.Wrapf(err, "opening template %s", path)
		}
		name := filepath.Base(path)
		buf.Reset()
		buf.ReadFrom(f)
		err = parse(name, buf.String())
		if err != nil {
			return nil, errors.Wrapf(err, "parsing template %s", path)
		}
	}
	return nil, nil
}
