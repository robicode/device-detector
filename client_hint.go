package devicedetector

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/robicode/device-detector/util"
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
}

func NewClientHint(cache *Cache, headers http.Header) *ClientHint {
	if headers == nil {
		return nil
	}

	browsers := extractBrowserList(headers)
	c := ClientHint{
		_cache:            cache,
		_headers:          headers,
		_browserList:      newBrowserList(browsers...),
		_mobile:           headers.Get("Sec-CH-UA-Mobile"),
		_model:            headers.Get("Sec-CH-UA-Model"),
		_platform:         headers.Get("Sec-CH-UA-Platform"),
		_platform_version: headers.Get("Sec-CH-UA-Platform-Version"),
	}

	c._appName = c.extractAppName()

	return &c
}

func (c *ClientHint) Brands() map[string]string {
	output := make(map[string]string)

	for _, browser := range c._hintBrowsers {
		output[browser.Name] = browser.Version
	}

	return output
}

func (c *ClientHint) BrowserName() string {
	return c._browserList.Reject(func(browser hintBrowser) bool {
		return browser.Name == "Chromium"
	}).Last().Name
}

func (c *ClientHint) Platform() string {
	if c == nil {
		return ""
	}
	return c._platform
}

func (c *ClientHint) PlatformVersion() string {
	return c._platform_version
}

func (c *ClientHint) Mobile() string {
	return c._mobile
}

func (c *ClientHint) Model() string {
	return c._model
}

// OSVersion returns the operating system version according to client hints.
func (c *ClientHint) OSVersion() string {
	if c == nil {
		return ""
	}
	if c._platform == "Windows" {
		return c.windowsVersion()
	}
	return c._platform_version
}

// OSFamily returns the operating system family according to client hints.
func (c *ClientHint) OSFamily() string {
	for family, members := range osFamilies {
		if util.InStrArray(c.osShortNmae(), members) {
			return family
		}
	}
	return ""
}

// OSName returns the operating system name according to client hints.
func (c *ClientHint) OSName() string {
	if c.IsAndroidApp() {
		return "Android"
	}
	return c.Platform()
}

// windowsVersion extracts the windows version from the *ClientHint.
func (c *ClientHint) windowsVersion() string {
	if c._platform_version == "" {
		return ""
	}

	majorVersion, err := strconv.Atoi(strings.Split(c._platform_version, ".")[0])
	if err != nil {
		return ""
	}

	if majorVersion < 1 {
		return ""
	}

	if majorVersion < 11 {
		return "10"
	} else {
		return "11"
	}
}

// IsAndroidApp returns if we are an android app based on client hints.
func (c *ClientHint) IsAndroidApp() bool {
	return util.InStrArray(c.appNameFromHeaders(), []string{"com.hisense.odinbrowser", "com.seraphic.openinet.pre",
		"com.appssppa.idesktoppcbrowser"})
}

// Private functions

func (c *ClientHint) appNameFromHeaders() string {
	if c == nil {
		return ""
	}
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

func (c *ClientHint) osShortNmae() string {
	for short, long := range operatingSystems {
		if strings.ToLower(c.OSName()) == strings.ToLower(long) {
			return short
		}
	}
	return ""
}
