package devicedetector

import "embed"

type Cache struct {
	Bot    BotCache
	Device DeviceCache
	OS     OSCache
}

//go:embed regexes/bots.yml
//go:embed regexes/client/browser_engine.yml
//go:embed regexes/client/browsers.yml
//go:embed regexes/client/feed_readers.yml
//go:embed regexes/client/hints/apps.yml
//go:embed regexes/client/hints/browsers.yml
//go:embed regexes/client/libraries.yml
//go:embed regexes/client/mediaplayers.yml
//go:embed regexes/client/mobile_apps.yml
//go:embed regexes/client/pim.yml
//go:embed regexes/device/cameras.yml
//go:embed regexes/device/car_browsers.yml
//go:embed regexes/device/consoles.yml
//go:embed regexes/device/mobiles.yml
//go:embed regexes/device/notebooks.yml
//go:embed regexes/device/portable_media_player.yml
//go:embed regexes/device/shell_tv.yml
//go:embed regexes/device/televisions.yml
//go:embed regexes/oss.yml
//go:embed regexes/vendorfragments.yml
var embeddedData embed.FS

// List of all files in cache
var cacheFiles = []string{
	"bots.yml",
	"client/browser_engine.yml",
	"client/browsers.yml",
	"client/feed_readers.yml",
	"client/hints/apps.yml",
	"client/hints/browsers.yml",
	"client/libraries.yml",
	"client/mediaplayers.yml",
	"client/mobile_apps.ym.",
	"client/pim.yml",
	"device/cameras.yml",
	"device/car_browsers.yml",
	"device/consoles.yml",
	"device/mobiles.yml",
	"device/notebooks.yml",
	"device/portable_media_player.yml",
	"device/shell_tv.yml",
	"device/televisions.yml",
	"oss.yml",
	"vendorfragments.yml",
}

func NewEmbeddedCache() (*Cache, error) {
	botCache, err := NewEmbeddedBotCache()
	if err != nil {
		return nil, err
	}

	deviceCache, err := NewEmbeddedDeviceCache()
	if err != nil {
		return nil, err
	}

	osCache, err := NewEmbeddedOSCache()
	if err != nil {
		return nil, err
	}

	return &Cache{
		Bot:    botCache,
		Device: deviceCache,
		OS:     osCache,
	}, nil
}
