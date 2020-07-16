package config

import "strings"
import "github.com/spf13/cast"
import "github.com/google/go-github/v32/github"

type config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}

type Config struct {
	projectName string
	projectType string
	projectDir  string
	projectRepo *github.Repository
	// These are config values, i.e. the settings declared outside of the [params] section.
	// This is the map used for looking for configuration values (baseURL etc.).
	// Values in this map can also be fetched from the params map above.
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

// GetLocal gets a configuration value set on language level. It will
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
