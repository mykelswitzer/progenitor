package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	OrgName  string    `yaml:"orgname,omitempty"`
}

var sttngsFilenames = []string{"progenitor.yaml"}

// DefaultConfig creates a copy of the default config
func DefaultSettings() *Settings {
	return &Settings{
	}
}

// LoadSettingsFromDefaultLocations looks for a settings file in the current directory, 
// and all parent directories walking up the tree. The closest settings file will be returned.
func LoadSettingsFromDefaultLocations() (*Settings, error) {
	settingsFile, err := findSettingsFile()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(filepath.Dir(settingsFile))
	if err != nil {
		return nil, fmt.Errorf("unable to enter config dir: %w", err)
	}
	return LoadSettings(settingsFile)
}

var path2regex = strings.NewReplacer(
	`.`, `\.`,
	`*`, `.+`,
	`\`, `[\\/]`,
	`/`, `[\\/]`,
)

// LoadSettings reads the gqlgen.yml config file
func LoadSettings(filename string) (*Settings, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	return ReadSettings(bytes.NewReader(b))
}

func ReadSettings(settingsFile io.Reader) (*Settings, error) {
	settings := DefaultSettings()

	dec := yaml.NewDecoder(settingsFile)
	dec.KnownFields(true)

	if err := dec.Decode(settings); err != nil {
		return nil, fmt.Errorf("unable to parse settings: %w", err)
	}

	if err := CompleteSettings(settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// CompleteSettings fills in the values to a settings loaded from YAML.
func CompleteSettings(settings *Settings) error {
	// defaultDirectives := map[string]DirectiveSettings{
	// 	"skip":        {SkipRuntime: true},
	// 	"include":     {SkipRuntime: true},
	// 	"deprecated":  {SkipRuntime: true},
	// 	"specifiedBy": {SkipRuntime: true},
	// }

	// for key, value := range defaultDirectives {
	// 	if _, defined := config.Directives[key]; !defined {
	// 		config.Directives[key] = value
	// 	}
	// }

	return nil
}


type StringList []string

func (a *StringList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single string
	err := unmarshal(&single)
	if err == nil {
		*a = []string{single}
		return nil
	}

	var multi []string
	err = unmarshal(&multi)
	if err != nil {
		return err
	}

	*a = multi
	return nil
}

func (a StringList) Has(file string) bool {
	for _, existing := range a {
		if existing == file {
			return true
		}
	}
	return false
}

func inStrSlice(haystack []string, needle string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}

	return false
}

// findSettingsFile searches for the config file in this directory and all parents up the tree
// looking for the closest match
func findSettingsFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get working dir to findSettingsFile: %w", err)
	}

	cfg := findSettingsFileInDir(dir)

	for cfg == "" && dir != filepath.Dir(dir) {
		dir = filepath.Dir(dir)
		cfg = findSettingsFileInDir(dir)
	}

	if cfg == "" {
		return "", os.ErrNotExist
	}

	return cfg, nil
}

func findSettingsFileInDir(dir string) string {
	for _, cfgName := range sttngsFilenames {
		path := filepath.Join(dir, cfgName)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func abs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return filepath.ToSlash(absPath)
}

type DirectiveSettings struct {
	SkipRuntime bool `yaml:"skip_runtime"`
}
