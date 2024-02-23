package config

type TerraformSettings struct {
	Workspaces StringList `yaml:"workspaces,omitempty"`
}

func (t *TerraformSettings) IsDefined() bool {
	return len(t.Workspaces) > 0
}
