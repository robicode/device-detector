package devicedetector

import "github.com/gijsbers/go-pcre"

type CachedBotProducer struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type CachedBot struct {
	Regex         string            `yaml:"regex"`
	compiledRegex pcre.Regexp       `yaml:"-"`
	compileError  error             `yaml:"-"`
	compiled      bool              `yaml:"-"`
	Name          string            `yaml:"name"`
	Category      string            `yaml:"category"`
	URL           string            `yaml:"url"`
	Producer      CachedBotProducer `yaml:"producer"`
}
