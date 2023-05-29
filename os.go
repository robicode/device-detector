package devicedetector

import (
	"strings"

	"github.com/robicode/device-detector/extractor"
	"github.com/robicode/device-detector/util"
	"github.com/robicode/device-detector/versionextractor"
)

type OS struct {
	_cache     *Cache
	_os        *CachedOS
	_userAgent string
}

// OS short codes mapped to long names
var operatingSystems = map[string]string{
	"AIX": "AIX",
	"AND": "Android",
	"ADR": "Android TV",
	"AMZ": "Amazon Linux",
	"AMG": "AmigaOS",
	"ATV": "tvOS",
	"ARL": "Arch Linux",
	"BTR": "BackTrack",
	"SBA": "Bada",
	"BEO": "BeOS",
	"BLB": "BlackBerry OS",
	"QNX": "BlackBerry Tablet OS",
	"BOS": "Bliss OS",
	"BMP": "Brew",
	"CAI": "Caixa Mágica",
	"CES": "CentOS",
	"CST": "CentOS Stream",
	"CLR": "ClearOS Mobile",
	"COS": "Chrome OS",
	"CRS": "Chromium OS",
	"CHN": "China OS",
	"CYN": "CyanogenMod",
	"DEB": "Debian",
	"DEE": "Deepin",
	"DFB": "DragonFly",
	"DVK": "DVKBuntu",
	"FED": "Fedora",
	"FEN": "Fenix",
	"FOS": "Firefox OS",
	"FIR": "Fire OS",
	"FOR": "Foresight Linux",
	"FRE": "Freebox",
	"BSD": "FreeBSD",
	"FYD": "FydeOS",
	"FUC": "Fuchsia",
	"GNT": "Gentoo",
	"GRI": "GridOS",
	"GTV": "Google TV",
	"HPX": "HP-UX",
	"HAI": "Haiku OS",
	"IPA": "iPadOS",
	"HAR": "HarmonyOS",
	"HAS": "HasCodingOS",
	"IRI": "IRIX",
	"INF": "Inferno",
	"JME": "Java ME",
	"KOS": "KaiOS",
	"KAN": "Kanotix",
	"KNO": "Knoppix",
	"KTV": "KreaTV",
	"KBT": "Kubuntu",
	"LIN": "GNU/Linux",
	"LND": "LindowsOS",
	"LNS": "Linspire",
	"LEN": "Lineage OS",
	"LBT": "Lubuntu",
	"LOS": "Lumin OS",
	"VLN": "VectorLinux",
	"MAC": "Mac",
	"MAE": "Maemo",
	"MAG": "Mageia",
	"MDR": "Mandriva",
	"SMG": "MeeGo",
	"MCD": "MocorDroid",
	"MON": "moonOS",
	"MIN": "Mint",
	"MLD": "MildWild",
	"MOR": "MorphOS",
	"NBS": "NetBSD",
	"MTK": "MTK / Nucleus",
	"MRE": "MRE",
	"WII": "Nintendo",
	"NDS": "Nintendo Mobile",
	"NOV": "Nova",
	"OS2": "OS/2",
	"T64": "OSF1",
	"OBS": "OpenBSD",
	"OWR": "OpenWrt",
	"OTV": "Opera TV",
	"ORD": "Ordissimo",
	"PAR": "Pardus",
	"PCL": "PCLinuxOS",
	"PLA": "Plasma Mobile",
	"PSP": "PlayStation Portable",
	"PS3": "PlayStation",
	"PUR": "PureOS",
	"RHT": "Red Hat",
	"REV": "Revenge OS",
	"ROS": "RISC OS",
	"ROK": "Roku OS",
	"RSO": "Rosa",
	"ROU": "RouterOS",
	"REM": "Remix OS",
	"RRS": "Resurrection Remix OS",
	"REX": "REX",
	"RZD": "RazoDroiD",
	"SAB": "Sabayon",
	"SSE": "SUSE",
	"SAF": "Sailfish OS",
	"SEE": "SeewoOS",
	"SIR": "Sirin OS",
	"SLW": "Slackware",
	"SOS": "Solaris",
	"SYL": "Syllable",
	"SYM": "Symbian",
	"SYS": "Symbian OS",
	"S40": "Symbian OS Series 40",
	"S60": "Symbian OS Series 60",
	"SY3": "Symbian^3",
	"TEN": "TencentOS",
	"TDX": "ThreadX",
	"TIZ": "Tizen",
	"TOS": "TmaxOS",
	"UBT": "Ubuntu",
	"WAS": "watchOS",
	"WTV": "WebTV",
	"WHS": "Whale OS",
	"WIN": "Windows",
	"WCE": "Windows CE",
	"WIO": "Windows IoT",
	"WMO": "Windows Mobile",
	"WPH": "Windows Phone",
	"WRT": "Windows RT",
	"XBX": "Xbox",
	"XBT": "Xubuntu",
	"YNS": "YunOS",
	"ZEN": "Zenwalk",
	"ZOR": "ZorinOS",
	"IOS": "iOS",
	"POS": "palmOS",
	"WOS": "webOS",
}

// List of operating system families.
var osFamilies = map[string][]string{
	"Android": {"AND", "CYN", "FIR", "REM", "RZD", "MLD", "MCD", "YNS", "GRI", "HAR",
		"ADR", "CLR", "BOS", "REV", "LEN", "SIR", "RRS"},
	"AmigaOS":        {"AMG", "MOR"},
	"BlackBerry":     {"BLB", "QNX"},
	"Brew":           {"BMP"},
	"BeOS":           {"BEO", "HAI"},
	"Chrome OS":      {"COS", "CRS", "FYD", "SEE"},
	"Firefox OS":     {"FOS", "KOS"},
	"Gaming Console": {"WII", "PS3"},
	"Google TV":      {"GTV"},
	"IBM":            {"OS2"},
	"iOS":            {"IOS", "ATV", "WAS", "IPA"},
	"RISC OS":        {"ROS"},
	"GNU/Linux": {
		"LIN", "ARL", "DEB", "KNO", "MIN", "UBT", "KBT", "XBT", "LBT", "FED",
		"RHT", "VLN", "MDR", "GNT", "SAB", "SLW", "SSE", "CES", "BTR", "SAF",
		"ORD", "TOS", "RSO", "DEE", "FRE", "MAG", "FEN", "CAI", "PCL", "HAS",
		"LOS", "DVK", "ROK", "OWR", "OTV", "KTV", "PUR", "PLA", "FUC", "PAR",
		"FOR", "MON", "KAN", "ZEN", "LND", "LNS", "CHN", "AMZ", "TEN", "CST",
		"NOV", "ROU", "ZOR",
	},
	"Mac":                   {"MAC"},
	"Mobile Gaming Console": {"PSP", "NDS", "XBX"},
	"Real-time OS":          {"MTK", "TDX", "MRE", "JME", "REX"},
	"Other Mobile":          {"WOS", "POS", "SBA", "TIZ", "SMG", "MAE"},
	"Symbian":               {"SYM", "SYS", "SY3", "S60", "S40"},
	"Unix":                  {"SOS", "AIX", "HPX", "BSD", "NBS", "OBS", "DFB", "SYL", "IRI", "T64", "INF"},
	"WebTV":                 {"WTV"},
	"Windows":               {"WIN"},
	"Windows Mobile":        {"WPH", "WMO", "WCE", "WRT", "WIO"},
	"Other Smart TV":        {"WHS"},
}

// List of desktop operating systems.
var desktopOSs = []string{"AmigaOS", "IBM", "GNU/Linux", "Mac", "Unix", "Windows", "BeOS", "Chrome OS"}

// NewOS returns a new *OS and tries to detect the OS of the given userAgent.
func NewOS(cache *Cache, userAgent string) *OS {
	os := cache.OS.Find(userAgent)

	return &OS{
		_cache:     cache,
		_os:        os,
		_userAgent: userAgent,
	}
}

// Short returns the short code for the OS.
func (o *OS) Short() string {
	if o._os == nil {
		return ""
	}

	for short, long := range operatingSystems {
		if strings.ToLower(o.Name()) == strings.ToLower(long) {
			return short
		}
	}
	return "UNK"
}

// Name returns the full name of the OS.
func (o *OS) Name() string {
	if o._os == nil {
		return ""
	}
	return extractor.New(o._userAgent, o._os.Regex, o._os.Name).Call()
}

// ShortName is an alias for Short(); it returns the short code for the OS.
func (o *OS) ShortName() string {
	return o.Short()
}

// Family returns the OS family
func (o *OS) Family() string {
	return o.familyToOS()
}

func (o *OS) familyToOS() string {
	if o._os == nil {
		return ""
	}

	for family, members := range osFamilies {
		if util.InStrArray(o.Short(), members) {
			return family
		}
	}
	return ""
}

// IsDesktop returns a best guess whether this OS is primarily a desktop OS.
func (o *OS) IsDesktop() bool {
	return util.InStrArray(o.Family(), desktopOSs)
}

// FullVersion returns the full version string of the OS. Not implemented yet.
func (o *OS) FullVersion() string {
	if o._os == nil {
		return ""
	}
	var versions []versionextractor.Version

	for _, version := range o._os.Versions {
		v := versionextractor.Version{
			Regex:   version.Regex,
			Version: version.Version,
		}
		versions = append(versions, v)
	}

	return versionextractor.New(o._userAgent, o._os.Regex, o._os.Version, versions).Call()
}
