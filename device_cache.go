package devicedetector

type DeviceCache interface {
	Delete(item string) DeviceCache
	Get(list string) *CachedDeviceList
	RegexesForHbbTV() *CacheFileList
	RegexesForShellTV() *CacheFileList
	RegexesForOthers() *CacheFileList
	regexFind(userAgent string, caches *CacheFileList) *CachedDevice
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
