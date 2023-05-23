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

// IsBot returns true if the userAgent given to NewBot represents a bot.
func (b *Bot) IsBot() bool {
	return b._bot != nil
}

func (b *Bot) Name() string {
	if b._bot != nil {
		return b._bot.Name
	}
	return ""
}

func (b *Bot) Category() string {
	if b._bot != nil {
		return b._bot.Category
	}
	return ""
}
