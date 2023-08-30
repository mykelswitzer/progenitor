package scaffold

import (
	"github.com/mykelswitzer/progenitor/internal/scaffold"
	"github.com/mykelswitzer/progenitor/pkg/config"
	"github.com/mykelswitzer/progenitor/pkg/prompt"
	"github.com/spf13/afero"
)

// ScaffoldDS is an interface that each template system
// implements to serve as the datasource for populating
// the scaffold and running the commands against
type ScaffoldDS interface {
	Init(*config.Config, []string, map[string]func(*Scaffold) error)
	GetName() string
	GetDescription() string
	GetPrompts() []prompt.PromptFunc
	GetSkipTemplates(*config.Config) []string
	GetProcessHooks(*config.Config) map[string]func(*Scaffold) error
	Populate(*string, afero.Fs) error
}

// Scaffolds is a keyed map of ScaffoldDS interfaces
type Scaffolds []ScaffoldDS

// Scaffold is a common struct that all template systems
// rely on to run the various scaffolding commands
type Scaffold = scaffold.Scaffold
