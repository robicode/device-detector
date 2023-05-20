package devicedetector

import "github.com/gijsbers/go-pcre"

type CachedModel struct {
	Regex         string      `yaml:"regex"`
	compiledRegex pcre.Regexp `yaml:"-"`
	compileError  error       `yaml:"-"`
	compiled      bool        `yaml:"-"`
	Brand         string      `yaml:"brand"`
	Name          string      `yaml:"model"`
	Type          string      `yaml:"device"`
}

type CachedDevice struct {
	Regex         string        `yaml:"regex"`
	Captures      []string      `yaml:"-"`
	compiledRegex pcre.Regexp   `yaml:"-"`
	compileError  error         `yaml:"-"`
	compiled      bool          `yaml:"-"`
	Type          string        `yaml:"device"`
	Name          string        `yaml:"model"`
	Models        []CachedModel `yaml:"models"`
	Brand         string        `yaml:"-"`
}
