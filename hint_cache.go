package devicedetector

import "log"

type HintCache interface {
	Find(id string) string
}

type EmbeddedHintCache struct {
	list map[string]string
}

func NewEmbeddedHintCache() (*EmbeddedHintCache, error) {
	files := NewCacheFileList("client/hints/browsers.yml", "client/hints/apps.yml")
	maps, err := parseHints(files)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &EmbeddedHintCache{
		list: maps,
	}, err
}

func (e *EmbeddedHintCache) Find(id string) string {
	if value, ok := e.list[id]; ok {
		return value
	}
	return ""
}
