package devicedetector

import (
	"log"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
)

type BotCache interface {
	Find(userAgent string) *CachedBot
}

type EmbeddedBotCache struct {
	bots []CachedBot
}

var botFiles = []string{
	"bots.yml",
}

func NewEmbeddedBotCache() (*EmbeddedBotCache, error) {
	files := NewCacheFileList(botFiles...)

	bots, err := parseBots(files)
	if err != nil {
		return nil, err
	}
	return &EmbeddedBotCache{
		bots: bots,
	}, nil
}

func (b *EmbeddedBotCache) Find(userAgent string) *CachedBot {
	var matches []CachedBot

	for _, bot := range b.bots {
		if !bot.compiled && bot.compileError == nil {
			re, err := pcre.Compile(util.FixupRegex(bot.Regex), pcre.CASELESS)
			if err != nil {
				bot.compileError = err
				log.Println(err)
				continue
			}
			bot.compiled = true
			bot.compiledRegex = re
		}

		if bot.compileError == nil {
			matcher := bot.compiledRegex.MatcherString(userAgent, 0)
			if matcher.Matches() {
				matches = append(matches, bot)
			}
		}
	}

	if len(matches) > 0 {
		return &matches[len(matches)-1]
	}
	return nil
}
