package devicedetector

import (
	"log"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
)

type ClientCache interface {
	Find(userAgent string) *CachedClient
}

type Engine struct {
	Default  string
	Versions map[string]string
}

type CachedClient struct {
	Regex         string
	compiledRegex pcre.Regexp
	compileError  error
	compiled      bool
	Name          string
	Path          string
	Type          string
	Version       string
	Engine        Engine
	URL           string
}

type EmbeddedClientCache struct {
	cleints []CachedClient
}

func NewEmbeddedClientCache() (*EmbeddedClientCache, error) {
	files := NewCacheFileList(clientFilenames...)
	clients, err := parseClients(files)
	if err != nil {
		return nil, err
	}

	return &EmbeddedClientCache{
		cleints: clients,
	}, nil
}

func (e *EmbeddedClientCache) Find(userAgent string) *CachedClient {
	reversed := util.ReverseArray(e.cleints)
	var matches []CachedClient

	for _, client := range reversed {
		if !client.compiled && client.compileError == nil {
			re, err := pcre.Compile(util.FixupRegex(client.Regex), pcre.CASELESS)
			if err != nil {
				client.compileError = err
				log.Println(err)
				continue
			}
			client.compiledRegex = re
			client.compiled = true
		}

		if client.compileError == nil {
			matcher := client.compiledRegex.MatcherString(userAgent, 0)
			if matcher.Matches() {
				matches = append(matches, client)
			}
		}
	}

	if len(matches) > 0 {
		return &matches[len(matches)-1]
	}
	return nil
}
