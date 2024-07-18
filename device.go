package devicedetector

import (
	"log"
	"regexp"
	"strings"

	"github.com/robicode/device-detector/extractor"
)

type Device struct {
	cache *Cache

	_device    *CachedDevice
	_model     *CachedModel
	_userAgent string
}

// NewDevice creates a new *Device.
func NewDevice(cache *Cache, userAgent string) *Device {
	device := &Device{
		cache:      cache,
		_device:    nil,
		_userAgent: userAgent,
	}

	entry := device.matchingRegex()
	if entry == nil {
		return device
	}

	device._device = entry
	return device
}

// Name returns the device name.
func (d *Device) Name() string {
	if d._device != nil {
		return d._device.Name
	}
	return ""
}

// Type returns the device type (desktop, smartphone...)
func (d *Device) Type() string {
	if d.isHbbTV() || d.isShellTV() {
		return "tv"
	}

	if d._device != nil {
		return d._device.Type
	}
	return ""
}

// Brand returns the device brand
func (d *Device) Brand() string {
	if d._device != nil {
		return d._device.Brand
	}
	return ""
}

func (d *Device) isHbbTV() bool {
	return regexp.MustCompile(`HbbTV/([1-9]{1}(?:\.[0-9]{1}){1,2})`).MatchString(d._userAgent)
}

func (d *Device) isShellTV() bool {
	return regexp.MustCompile(`[a-z]+[ _]Shell[ _]\w{6}`).MatchString(d._userAgent)
}

func (d *Device) matchingRegex() *CachedDevice {
	var regexList = NewCacheFileList()

	if d.isHbbTV() {
		regexList = d.cache.Device.RegexesForHbbTV()
	} else {
		regexList = d.cache.Device.RegexesForOthers()
	}

	if d.isShellTV() {
		regexList = d.cache.Device.RegexesForShellTV()
	}

	if regexList == nil {
		log.Println("BUG: regexList is nil! This should not happen!")
		return nil
	}

	device := d.cache.Device.regexFind(d._userAgent, regexList)
	if device == nil {
		return nil
	}

	if len(device.Models) > 0 {
		model := device.FindModel(d._userAgent)
		if model != nil {
			d._model = model
			if strings.TrimSpace(model.Brand) != "" {
				device.Brand = model.Brand
			}
			if strings.TrimSpace(model.Name) != "" {
				name := extractor.New(d._userAgent, model.Regex, model.Name).Call()
				device.Name = name
			}
			if strings.TrimSpace(model.Type) != "" {
				device.Type = model.Type
			}
		}
	} else {
		name := extractor.New(d._userAgent, device.Regex, device.Name).Call()
		device.Name = name
	}
	return device
}

func (d *Device) IsKnown() bool {
	return d._device != nil
}
