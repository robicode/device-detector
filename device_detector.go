package devicedetector

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
	"github.com/robicode/device-detector/versionextractor"
	"github.com/robicode/version"
)

type DeviceDetector struct {
	_bot       *Bot
	_cache     *Cache
	_client    *Client
	_hint      *ClientHint
	_device    *Device
	_os        *OS
	_userAgent string
}

// New returns a newly initialised *DeviceDetector.
func New(cache *Cache, userAgent string, headers ...http.Header) *DeviceDetector {
	var h http.Header
	var hint *ClientHint
	if headers != nil {
		h = headers[0]
	}

	if headers != nil && len(h) != 0 {
		hint = NewClientHint(cache, h)
	}

	bot := NewBot(cache, userAgent)
	client := NewClient(cache, userAgent)

	device := NewDevice(cache, userAgent)
	os := NewOS(cache, userAgent)

	return &DeviceDetector{
		_bot:       bot,
		_cache:     cache,
		_client:    client,
		_hint:      hint,
		_device:    device,
		_os:        os,
		_userAgent: userAgent,
	}
}

// BotName returns the bot name or an empty string if client is not a bot.
func (d *DeviceDetector) BotName() string {
	if d._bot == nil {
		return ""
	}
	return d._bot.Name()
}

// IsBot returns true if a bot is detected.
func (d *DeviceDetector) IsBot() bool {
	if d._bot == nil {
		return false
	}
	return d._bot.IsBot()
}

// IsKnown returns true if this client is known.
func (d *DeviceDetector) IsKnown() bool {
	if d._client == nil {
		return false
	}
	return d._client.IsKnown()
}

// OSFullVersion returns the operating system version.
func (d *DeviceDetector) OSFullVersion() string {
	if d.skipOSVersion() {
		return ""
	}
	if d._hint != nil {
		if d._hint.OSVersion() != "" {
			return d._hint.OSVersion()
		}
	}
	return d._os.FullVersion()
}

// Name returns the client name.
func (d *DeviceDetector) Name() string {
	if d.needsMobileFix() {
		return d._client.Name()
	}
	if d._hint != nil && d._hint.BrowserName() != "" {
		return d._hint.BrowserName()
	}
	return d._client.Name()
}

// FullVersion returns the client version.
func (d *DeviceDetector) FullVersion() string {
	if d._hint != nil {
		return d._hint.PlatformVersion()
	}
	return d._client.FullVersion()
}

// OSFamily returns the operating system family.
func (d *DeviceDetector) OSFamily() string {
	if d.needsLinuxFix() {
		return "GNU/Linux"
	}
	// client_hint.os_family || os.family || client_hint.platform
	if d._hint != nil {
		if d._hint.OSFamily() != "" {
			return d._hint.OSFamily()
		}
	}
	if d._os != nil {
		if d._os.Family() != "" {
			return d._os.Family()
		}
	}
	if d._hint != nil {
		if d._hint.OSFamily() != "" {
			return d._hint.Platform()
		}
	}
	return ""
}

// OSName returns the operating system name.
func (d *DeviceDetector) OSName() string {
	if d.needsLinuxFix() {
		return "GNU/Linux"
	}
	if d._hint != nil && d._hint.OSName() != "" {
		return d._hint.OSName()
	}
	if d._os != nil && d._os.Name() != "" {
		return d._os.Name()
	}
	if d._hint != nil {
		return d._hint.Platform()
	}
	return ""
}

// DeviceName returns the detected device name.
func (d *DeviceDetector) DeviceName() string {
	if d._device != nil && d._device.Name() != "" {
		return d._device.Name()
	}

	if d._hint != nil && d._hint.Model() != "" {
		return d._hint.Model()
	}
	return d.fixForXMusic()
}

// DeviceBrand returns the detected device brand.
func (d *DeviceDetector) DeviceBrand() string {
	if d._device == nil {
		return ""
	}
	// Assume all devices running iOS / Mac OS are from Apple
	brand := d._device.Brand()
	if util.InStrArray(d.OSName(), []string{"Apple TV", "iOS", "Mac"}) {
		brand = "Apple"
	}
	return brand
}

// DeviceType attempts to detect the device type.
func (d *DeviceDetector) DeviceType() string {
	t := d._device.Type()

	// Chrome on Android passes the device type based on the keyword 'Mobile'
	// If it is present the device should be a smartphone, otherwise it's a tablet
	// See https://developer.chrome.com/multidevice/user-agent#chrome_for_android_user_agent
	// Note: We do not check for browser (family) here, as there might be mobile apps using Chrome,
	// that won't have a detected browser, but can still be detected. So we check the useragent for
	// Chrome instead.
	if t == "" && d.OSFamily() == "Android" && pcre.MustCompile(buildRegex(`Chrome\/[\.0-9]*`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches() {
		if pcre.MustCompile(buildRegex(`(?:Mobile|eliboM) Safari\/`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches() {
			t = "smartphone"
		} else if pcre.MustCompile(buildRegex(`(?!Mobile )Safari\/`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches() {
			t = "tablet"
		}
	}

	// Some UA contain the fragment 'Android; Tablet;' or 'Opera Tablet', so we assume those devices
	// # as tablets
	if t == "" && d.hasAndroidTabletFragment() || d.isOperaTablet() {
		t = "tablet"
	}

	// Some user agents simply contain the fragment 'Android; Mobile;', so we assume those devices
	// as smartphones
	if t == "" && d.hasAndroidMobileFragment() {
		t = "smartphone"
	}

	// Android up to 3.0 was designed for smartphones only. But as 3.0,
	// which was tablet only, was published too late, there were a
	// bunch of tablets running with 2.x With 4.0 the two trees were
	// merged and it is for smartphones and tablets
	//
	// So were are expecting that all devices running Android < 2 are
	// smartphones Devices running Android 3.X are tablets. Device type
	// of Android 2.X and 4.X+ are unknown
	if t == "" && d.OSName() == "Android" && d._os.FullVersion() != "" {
		fullVersion := version.New2(d._os.FullVersion())
		if fullVersion.Compare(versionextractor.MajorVersion2) == 1 {
			t = "smartphone"
		} else if fullVersion.Compare(versionextractor.MajorVersion3) <= 0 &&
			fullVersion.Compare(versionextractor.MajorVersion4) == 1 {
			t = "tablet"
		}
	}

	// All detected feature phones running android are more likely a smartphone
	if t == "feature phone" && d.OSFamily() == "Android" {
		t = "smartphone"
	}

	// All unknown devices under running Java ME are more likely a features phones
	if t == "" && d.OSName() == "Java ME" {
		t = "feature phone"
	}

	// According to http://msdn.microsoft.com/en-us/library/ie/hh920767(v=vs.85).aspx
	// Internet Explorer 10 introduces the "Touch" UA string token. If this token is present at the
	// end of the UA string, the computer has touch capability, and is running Windows 8 (or later).
	// This UA string will be transmitted on a touch-enabled system running Windows 8 (RT)
	//
	// As most touch enabled devices are tablets and only a smaller part are desktops/notebooks we
	// assume that all Windows 8 touch devices are tablets.
	if t == "" && d.isTouchEnabled() && (d.OSName() == "Windows RT" || (d.OSName() == "Windows" && d.OSFullVersion() != "" &&
		version.New2(d.OSFullVersion()).Compare(versionextractor.MajorVersion8) <= 0)) {
		t = "tablet"
	}

	// All devices running Opera TV Store are assumed to be a tv
	if d.isOperaTVStore() {
		t = "tv"
	}

	// All devices that contain Andr0id in string are assumed to be a tv
	if pcre.MustCompile(buildRegex(`Andr0id|Android TV`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches() {
		t = "tv"
	}

	// All devices running Tizen TV or SmartTV are assumed to be a tv
	if d.isTizenTV() {
		t = "tv"
	}

	// Devices running Kylo or Espital TV Browsers are assumed to be a TV
	if t == "" && (strings.Contains(d.Name(), "Kylo") || strings.Contains(d.Name(), "Espial TV Browser")) {
		t = "tv"
	}

	// All devices containing TV fragment are assumed to be a tv
	if t == "" && pcre.MustCompile(buildRegex(`\(TV;`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches() {
		t = "tv"
	}

	var hasDesktop bool
	if t != "desktop" && d.hasDesktopString() && d.hasDesktopFragment() {
		hasDesktop = true
	}

	if hasDesktop {
		t = "desktop"
	}

	// set device type to desktop for all devices running a desktop os that were not detected as
	// # another device type
	if t != "" || !d.isDesktop() {
		return t
	}
	return "desktop"
}

// Private functions

// https://github.com/matomo-org/device-detector/blob/be1c9ef486c247dc4886668da5ed0b1c49d90ba8/Parser/Client/Browser.php#L772
// Fix mobile browser names e.g. Chrome => Chrome Mobile
func (d *DeviceDetector) needsMobileFix() bool {
	if d._hint == nil || d._client == nil {
		return false
	}

	return d._client.Name() == fmt.Sprintf("%s Mobile", d._hint.BrowserName())
}

func (d *DeviceDetector) needsLinuxFix() bool {
	if d._hint == nil || d._os == nil {
		return false
	}
	if d._hint.Platform() == "Linux" && d._os.Name() == "Android" && d._hint._mobile == "?0" {
		return true
	}
	return false
}

func (d *DeviceDetector) fixForXMusic() string {
	if strings.Contains(d._userAgent, "X-music Ⅲ") {
		return "X-Music III"
	}
	return ""
}

func (d *DeviceDetector) skipOSVersion() bool {
	return d._hint.OSFamily() != "" && d._hint.OSFamily() != d._os.Family()
}

func (d *DeviceDetector) hasAndroidTabletFragment() bool {
	return pcre.MustCompile(buildRegex(`Android( [\.0-9]+)?; Tablet;`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) hasAndroidMobileFragment() bool {
	return pcre.MustCompile(buildRegex(`Android( [\.0-9]+)?; Mobile;`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) hasDesktopFragment() bool {
	return pcre.MustCompile(buildRegex(`Desktop (x(?:32|64)|WOW64);`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) isOperaTablet() bool {
	return pcre.MustCompile(buildRegex(`Opera Tablet`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) isTouchEnabled() bool {
	return pcre.MustCompile(buildRegex(`Touch`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) isOperaTVStore() bool {
	return pcre.MustCompile(buildRegex(`Opera TV Store`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) isTizenTV() bool {
	return pcre.MustCompile(buildRegex(`SmartTV|Tizen.+ TV .+$`), pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) hasDesktopString() bool {
	return pcre.MustCompile(`Desktop`, pcre.CASELESS).MatcherString(d._userAgent, 0).Matches()
}

func (d *DeviceDetector) isDesktop() bool {
	if d.OSName() == "" || d.OSName() == "UNK" {
		return false
	}

	// Check for browsers available for mobile devices only
	if d.usesMobileBrowser() {
		return false
	}

	return d._os.IsDesktop()
}

func (d *DeviceDetector) usesMobileBrowser() bool {
	return d._client.IsBrowser() && d._client.isMobileOnlyBrowser()
}

func buildRegex(src string) string {
	return fmt.Sprintf(`(?:^|[^A-Z0-9\_\-])(?:%s)`, src)
}
