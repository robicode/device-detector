package devicedetector

import (
	"strings"

	"github.com/robicode/device-detector/versionextractor"
)

type Client struct {
	_cache     *Cache
	_client    *CachedClient
	_userAgent string
}

var clientFilenames = []string{
	"client/feed_readers.yml",
	"client/mobile_apps.yml",
	"client/mediaplayers.yml",
	"client/pim.yml",
	"client/browsers.yml",
	"client/libraries.yml",
}

func NewClient(cache *Cache, userAgent string) *Client {
	c := Client{
		_cache:     cache,
		_userAgent: userAgent,
	}
	c._client = cache.Client.Find(userAgent)

	return &c
}

func (c *Client) Name() string {
	if c != nil {
		if c._client != nil {
			return c._client.Name
		}
	}
	return ""
}

func (c *Client) FullVersion() string {
	if c != nil {
		if c._client != nil {
			return versionextractor.New(c._userAgent, c._client.Regex, c._client.Version, nil).Call()
		}
	}
	return ""
}

func (c *Client) isMobileOnlyBrowser() bool {
	return mobileOnlyBrowser(c.Name())
}

func (c *Client) IsKnown() bool {
	return c._client != nil
}

func (c *Client) IsBrowser() bool {
	return strings.Contains(c._client.Path, "client/browsers.yml")
}
