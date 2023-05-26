package devicedetector

import (
	"log"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
)

type CachedOSVersion struct {
	Regex   string
	Version string
}

type CachedOS struct {
	Regex         string
	compileError  error
	compiledRegex pcre.Regexp
	compiled      bool
	Name          string
	Version       string
	Versions      []CachedOSVersion
}

type OSCache interface {
	Find(userAgent string) *CachedOS
}

type EmbeddedOSCache struct {
	osList []CachedOS
}

var osFiles = []string{
	"oss.yml",
}

func NewEmbeddedOSCache() (*EmbeddedOSCache, error) {
	files := NewCacheFileList(osFiles...)

	oss, err := parseOSs(files)
	if err != nil {
		return nil, err
	}

	return &EmbeddedOSCache{
		osList: oss,
	}, nil
}

func (e *EmbeddedOSCache) Find(userAgent string) *CachedOS {
	for _, os := range e.osList {
		if !os.compiled && os.compileError == nil {
			re, err := pcre.Compile(util.FixupRegex(os.Regex), pcre.CASELESS)
			if err != nil {
				os.compileError = err
				log.Println(err)
				continue
			}
			os.compiled = true
			os.compiledRegex = re
		}

		if os.compileError == nil {
			matcher := os.compiledRegex.MatcherString(userAgent, 0)
			if matcher.Matches() {
				return &os
			}
		}
	}
	return nil
}
