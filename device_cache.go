package devicedetector

import (
	"log"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
)

type DeviceCache interface {
	Delete(item string) DeviceCache
	Get(list string) *CachedDeviceList
	RegexesForHbbTV() *CacheFileList
	RegexesForShellTV() *CacheFileList
	RegexesForOthers() *CacheFileList
	regexFind(userAgent string, caches *CacheFileList) *CachedDevice
}

type EmbeddedDeviceCache struct {
	devices map[string]CachedDeviceList
}

// List of YAML files containing device specifications.
var deviceFilenames = []string{
	"device/televisions.yml",
	"device/shell_tv.yml",
	"device/notebooks.yml",
	"device/consoles.yml",
	"device/car_browsers.yml",
	"device/cameras.yml",
	"device/portable_media_player.yml",
	"device/mobiles.yml",
}

// NewEmbeddedDeviceCache creates a new embedded device cache, whereby the device
// tree is loaded from resources contained within the package/application binary.
func NewEmbeddedDeviceCache() (*EmbeddedDeviceCache, error) {
	files := NewCacheFileList(deviceFilenames...)
	devices, err := parse(files)
	if err != nil {
		return nil, err
	}

	return &EmbeddedDeviceCache{
		devices: devices,
	}, nil
}

func (e *EmbeddedDeviceCache) Delete(list string) DeviceCache {
	newDevicesList := make(map[string]CachedDeviceList)
	deletion := false
	for name, devices := range e.devices {
		if name == list {
			deletion = true
		} else {
			newDevicesList[name] = devices
		}
	}

	if deletion {
		newCache, err := NewEmbeddedDeviceCache()
		if err != nil {
			return nil
		}

		newCache.devices = newDevicesList
		return newCache
	}

	return e
}

// Get returns the devices with the associated list name.
func (e *EmbeddedDeviceCache) Get(list string) *CachedDeviceList {
	devices, ok := e.devices[list]
	if !ok {
		return nil
	}
	return &devices
}

func (e *EmbeddedDeviceCache) RegexesForHbbTV() *CacheFileList {
	return NewCacheFileList().Exclusive("regexes/device/televisions.yml")
}

func (e *EmbeddedDeviceCache) RegexesForShellTV() *CacheFileList {
	return NewCacheFileList().Exclusive("regexes/device/shell_tv.yml")
}

func (e *EmbeddedDeviceCache) RegexesForOthers() *CacheFileList {
	return NewCacheFileList().Exclude("regexes/device/televisions.yml")
}

func (e *EmbeddedDeviceCache) regexFind(userAgent string, caches *CacheFileList) *CachedDevice {
	for listfile, d := range e.devices {
		if !caches.Includes(listfile) {
			continue
		}

		for name, device := range d.list {
			if !device.compiled && device.compileError == nil {
				re, err := pcre.Compile(util.FixupRegex(device.Regex), pcre.CASELESS)
				if err != nil {
					device.compileError = err
					log.Println(err)
					continue
				}
				device.compiled = true
				device.compiledRegex = re
			}

			if device.compileError == nil {
				matcher := device.compiledRegex.MatcherString(userAgent, 0)
				if matcher.Matches() {
					device.Brand = name
					var captures []string
					for i := 0; i == matcher.Groups(); i++ {
						captures = append(captures, matcher.GroupString(i))
					}
					device.Captures = captures
					return &device
				}
			}
		}
	}
	return nil
}
