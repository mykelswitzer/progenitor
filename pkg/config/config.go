package config

import "strings"
import "github.com/spf13/cast"

// below are constants for the various keys
// held in the config...
// these should be added to when new prompts
// are setup to hold more info in the config
const CFG_DB_REQ   = "dbRequired"
const CFG_DB_MDL   = "dbModel"
const CFG_GQL_REQ  = "gqlRequired"

const CFG_ORG_NAME = "ghOrgName"
const CFG_PRJ_DIR  = "projectDir"
const CFG_PRJ_NAME = "projectName"
const CFG_PRJ_REPO = "projectRepo"
const CFG_PRJ_TEAM = "projectTeam"
const CFG_PRJ_TYPE = "projectType"

const CFG_RPT_REQ  = "reportingRequired"
const CFG_TF_RUN   = "runTerraform"

type config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}

type Config struct {
	settings map[string]interface{}
}

func New() *Config {
	return &Config{settings: make(map[string]interface{})}
}

// GetString returns the value associated with the key as a string.
func (c *Config) GetString(key string) string { return cast.ToString(c.Get(key)) }

// GetInt returns the value associated with the key as an int.
func (c *Config) GetInt(key string) int { return cast.ToInt(c.Get(key)) }

// GetBool returns the value associated with the key as a boolean.
func (c *Config) GetBool(key string) bool { return cast.ToBool(c.Get(key)) }

// Get gets a configuration value set on language level. It will
// not fall back to any global value.
// It will return nil if a value with the given key cannot be found.
func (c *Config) Get(key string) interface{} {
	if c == nil {
		panic("config not set")
	}
	key = strings.ToLower(key)
	if v, ok := c.settings[key]; ok {
		return v
	}

	return nil
}

// Set sets the value for the key in the config's params.
func (c *Config) Set(key string, value interface{}) {
	if c == nil {
		panic("config not set")
	}
	key = strings.ToLower(key)
	c.settings[key] = value
}

// IsSet checks whether the key is set in the language or the related config store.
func (c *Config) IsSet(key string) bool {
	key = strings.ToLower(key)
	if _, ok := c.settings[key]; ok {
		return true
	}
	return false
}

func (c *Config) GetSettings() map[string]interface{} {
	return c.settings
}
