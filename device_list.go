package devicedetector

import (
	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
)

// A DeviceList holds a list of devices loaded from YAML files.
type CachedDeviceList struct {
	list map[string]CachedDevice
}

// NewDeviceList creates a new list of devices.
func NewDeviceList() *CachedDeviceList {
	devices := make(map[string]CachedDevice)

	return &CachedDeviceList{
		list: devices,
	}
}

// Append returns a copy of the *CachedDeviceList with the given entry appended
// to the end of the list. If the new item already appears in the list, Append
// returns the original list unmodified.
func (d *CachedDeviceList) Append(name string, entry CachedDevice) *CachedDeviceList {
	for _, item := range d.list {
		if item.Regex == entry.Regex {
			return d
		}
	}

	dl := NewDeviceList()
	for key, value := range d.list {
		dl.list[key] = value
	}

	dl.list[name] = entry
	return dl
}

// Delete deletes an item from the *CachedDeviceList. If the item is not
// found, the original list is returned.
func (d *CachedDeviceList) Delete(item string) *CachedDeviceList {
	if _, ok := d.list[item]; ok {
		devices := NewDeviceList()
		for name, device := range d.list {
			if name == item {
				continue
			}

			devices.Append(name, device)
		}
		return devices
	}
	return d
}

// findModel attempts to locate a matching model, if any.
func (e *CachedDevice) FindModel(userAgent string) *CachedModel {
	if len(e.Models) == 0 {
		return nil
	}

	for _, model := range e.Models {
		if !model.compiled && model.compileError == nil {
			re, err := pcre.Compile(util.FixupRegex(model.Regex), pcre.CASELESS)
			if err != nil {
				model.compileError = err
				continue
			}
			model.compiled = true
			model.compiledRegex = re
		}

		// pcre uses Regexp, not *Regexp, so we can't check it for nil.
		if model.compileError == nil {
			modelMatcher := model.compiledRegex.MatcherString(userAgent, 0)
			if modelMatcher.Matches() {
				return &model
			}
		}
	}
	return nil
}
