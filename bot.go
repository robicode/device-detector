package devicedetector

type Bot struct {
	_bot       *CachedBot
	_cache     *Cache
	_userAgent string
}

func NewBot(cache *Cache, userAgent string) *Bot {
	if cache == nil {
		return nil
	}

	bot := cache.Bot.Find(userAgent)

	return &Bot{
		_bot:       bot,
		_cache:     cache,
		_userAgent: userAgent,
	}
}
