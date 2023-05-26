package devicedetector

import (
	"fmt"
	"net/http"
	"strings"
)

type ClientHint struct {
	_appName          string
	_cache            *Cache
	_headers          http.Header
	_browserList      *BrowserList
	_hintBrowsers     []hintBrowser
	_mobile           string
	_model            string
	_platform         string
	_platform_version string
	_userAgent        string
}

func NewClientHint(cache *Cache, userAgent string, headers http.Header) *ClientHint {
	if headers == nil {
		return nil
	}

	browsers := extractBrowserList(headers)
	c := ClientHint{
		_cache:            cache,
		_headers:          headers,
		_userAgent:        userAgent,
		_browserList:      newBrowserList(browsers...),
		_mobile:           headers.Get("Sec-CH-UA-Mobile"),
		_model:            headers.Get("Sec-CH-UA-Model"),
		_platform:         headers.Get("Sec-CH-UA-Platform"),
		_platform_version: headers.Get("Sec-CH-UA-Platform-Version"),
	}

	c._appName = c.extractAppName()

	return &c
}

func (c *ClientHint) appNameFromHeaders() string {
	if c._headers == nil {
		return ""
	}

	if c._headers.Get("http-x-requested-with") != "" {
		return c._headers.Get("http-x-requested-with")
	}

	if c._headers.Get("X-Requested-With") != "" {
		return c._headers.Get("X-Requested-With")
	}

	return ""
}

func (c *ClientHint) extractAppName() string {
	requestedWith := c.appNameFromHeaders()
	if requestedWith == "" {
		return ""
	}
	return c._cache.Hint.Find(requestedWith)
}

func extractBrowserList(headers http.Header) []hintBrowser {
	if headers.Get("Sec-CH-UA") == "" {
		return nil
	}

	var hintBrowsers []hintBrowser

	components := strings.Split(headers.Get("Sec-CH-UA"), ", ")
	for _, component := range components {
		componentAndVersion := strings.Split(strings.ReplaceAll(component, "\"", ""), ";v=")
		name := nameFromKnownBrowsers(componentAndVersion[0])
		if name == "" {
			continue
		}

		hb := hintBrowser{
			Name:    name,
			Version: componentAndVersion[1],
		}
		hintBrowsers = append(hintBrowsers, hb)
	}
	return hintBrowsers
}

// https://github.com/matomo-org/device-detector/blob/be1c9ef486c247dc4886668da5ed0b1c49d90ba8/Parser/Client/Browser.php#L865
func nameFromKnownBrowsers(name string) string {
	// https://github.com/matomo-org/device-detector/blob/be1c9ef486c247dc4886668da5ed0b1c49d90ba8/Parser/Client/Browser.php#L628
	if name == "Google Chrome" {
		return "Chrome"
	}

	for _, browser := range availableBrowsers {
		if browser == name {
			return browser
		}

		if strings.ReplaceAll(browser, " ", "") == strings.ReplaceAll(name, " ", "") {
			return browser
		}

		if browser == strings.ReplaceAll(name, "Browser", "") {
			return browser
		}

		if browser == strings.ReplaceAll(name, " Browser", "") {
			return browser
		}

		if browser == fmt.Sprintf("%s Browser", name) {
			return browser
		}
	}
	return ""
}

func (c *ClientHint) BrowserName() string {
	return c._browserList.Reject(func(browser hintBrowser) bool {
		return browser.Name == "Chromium"
	}).Last().Name
}

func (c *ClientHint) isIridium() bool {
	if len(c._browserList.list) == 0 {
		return false
	}

	browsers := c._browserList.Reject(func(browser hintBrowser) bool {
		if browser.Name == "Chromium" {
			if browser.Version == "2022" || browser.Version == "2022.04" {
				return true
			}
		}
		return false
	})
	return len(browsers.list) != 0
}
