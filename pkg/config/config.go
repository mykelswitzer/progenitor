package config

import "strings"
import "github.com/spf13/cast"

type config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}

type Config struct {
	settings *Settings
	inputs   map[string]interface{}
}

func New(settingsFile string) (*Config, error) {

	s, err := LoadSettings(settingsFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		settings: s,
		inputs:   make(map[string]interface{}),
	}

	return cfg, nil
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
	if v, ok := c.inputs[key]; ok {
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
	c.inputs[key] = value
}

// IsSet checks whether the key is set in the language or the related config store.
func (c *Config) IsSet(key string) bool {
	key = strings.ToLower(key)
	if _, ok := c.inputs[key]; ok {
		return true
	}
	return false
}

func (c *Config) GetInputs() map[string]interface{} {
	return c.inputs
}

func (c *Config) GetSettings() *Settings {
	return c.settings
}
