package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	GitHub   GitHubSettings `yaml:"github,omitempty"`
	Branches BranchSettings `yaml:"branches,omitempty"`
	Teams    StringList     `yaml:"teams,omitempty"`
}

type GitHubSettings struct {
	Organization string      `yaml:"organization,omitempty"`
	Creds		 GitHubCreds `yaml:"creds,omitempty"`
	App          GitHubApp   `yaml:"app,omitempty"`
}

func (s *GitHubSettings) IsDefined() bool {
	return s.Organization!="" && (s.Creds.IsDefined() || s.App.IsDefined())
}

func (s *GitHubSettings) UseCreds() bool {
	return s.Creds.IsDefined()
}

func (s *GitHubSettings) UseApp() bool {
	return s.App.IsDefined()
}

type GitHubApp struct {
	ID        	 int64 `yaml:"id,omitempty"`
	Key       	 string `yaml:"key,omitempty"`
	Installation int64 `yaml:"installation,omitempty"`
}

func (s *GitHubApp) IsDefined() bool {
	return s.ID!=0 || s.Key!="" || s.Installation!=0
}

type GitHubCreds struct {
	PAT        	 string `yaml:"pat,omitempty"`
}

func (s *GitHubCreds) IsDefined() bool {
	return s.PAT!=""
}

type BranchSettings struct {
	Default string `yaml:"default,omitempty"`
}

// LoadSettings reads the progenitor.yml settings file
func LoadSettings(settingsFile string) (*Settings, error) {
	b := []byte(settingsFile)
	return ReadSettings(bytes.NewReader(b))
}

func ReadSettings(settingsFile io.Reader) (*Settings, error) {
	s := &Settings{}

	dec := yaml.NewDecoder(settingsFile)
	dec.KnownFields(true)

	if err := dec.Decode(s); err != nil {
		return nil, fmt.Errorf("unable to parse settings: %w", err)
	}

	if err := s.check(); err != nil {
		return nil, err
	}

	if err := CompleteSettings(s); err != nil {
		return nil, err
	}

	return s, nil
}

// CompleteSettings fills in the values to a settings loaded from YAML.
func CompleteSettings(settings *Settings) error {
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

func (s *Settings) check() error {
	if !s.GitHub.IsDefined() {
		return errors.New("github settings are required")
	}

	return nil
}
