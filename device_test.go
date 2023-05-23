package devicedetector

import (
	"testing"
)

func Test_NewDevice(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Error("error creating device cache:", err)
		t.Fail()
		return
	}

	device := NewDevice(cache, "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12B466 [FBDV/iPhone7,2]")
	if device == nil {
		t.Error("expected device to != nil")
		t.Fail()
		return
	}
}

func TestDevice_Name(t *testing.T) {
	// Setup cache
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Error("error creating device cache:", err)
		t.Fail()
		return
	}

	// when models are nested

	device := NewDevice(cache, "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12B466 [FBDV/iPhone7,2]")
	if device == nil {
		t.Error("expected device to != nil")
		t.Fail()
		return
	}

	name := device.Name()
	if name != "iPhone 6" {
		t.Errorf("expected 'iPhone 6' but got '%s'", name)
		t.Fail()
		return
	}

	// when it cannot find a device name

	device = NewDevice(cache, "UNKNOWN MODEL NAME")
	name = device.Name()
	if name != "" {
		t.Errorf("expected name to be blank but got '%s'", name)
		t.Fail()
		return
	}

	// when models are NOT nested

	device = NewDevice(cache, "AIRNESS-AIR99/REV 2.2.1/Teleca Q03B1")

	name = device.Name()
	if name != "AIR99" {
		t.Errorf("expected 'AIR99' but got '%s'\nDevice: %v", name, device)
		t.Fail()
		return
	}
}

type DeviceTypeTest struct {
	Description   string
	UserAgent     string
	ExpectedBrand string
	ExpectedName  string
	ExpectedType  string
}

var typeTests = []DeviceTypeTest{
	{
		Description:   "Mobiles",
		UserAgent:     "Mozilla/5.0 (Linux; Android 4.4.2; es-us; SAMSUNG SM-G900F Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko)",
		ExpectedBrand: "Samsung",
		ExpectedName:  "Galaxy S5",
		ExpectedType:  "smartphone",
	},
	{
		Description:   "Cameras",
		UserAgent:     "Mozilla/5.0 (Linux; U; Android 4.0; xx-xx; EK-GC100 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		ExpectedBrand: "Samsung",
		ExpectedName:  "Galaxy Camera",
		ExpectedType:  "camera",
	},
	{
		Description:   "Car Browser",
		UserAgent:     "Mozilla/5.0 (X11; Linux) AppleWebKit/534.34 (KHTML, like Gecko) QtCarBrowser Safari/534.34",
		ExpectedBrand: "Tesla",
		ExpectedName:  "Model S",
		ExpectedType:  "car browser",
	},
	{
		Description:   "(Gaming) Consoles",
		UserAgent:     "Opera/9.30 (Nintendo Wii; U; ; 2047-7;en)",
		ExpectedBrand: "Nintendo",
		ExpectedName:  "Wii",
		ExpectedType:  "console",
	},
	{
		Description:   "Portable Media Players",
		UserAgent:     "Mozilla/5.0 (iPod touch; CPU iPhone OS 7_0_6 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B651 Safari/9537.53",
		ExpectedBrand: "Apple",
		ExpectedName:  "iPod Touch",
		ExpectedType:  "portable media player",
	},
	{
		Description:   "Televisions",
		UserAgent:     "Mozilla/5.0 (Unknown; Linux armv7l) AppleWebKit/537.1+ (KHTML, like Gecko) Safari/537.1+ HbbTV/1.1.1 ( ;LGE ;NetCast 4.0 ;03.10.81 ;1.0M ;)",
		ExpectedBrand: "LG",
		ExpectedName:  "NetCast 4.0",
		ExpectedType:  "tv",
	},
}

func TestDevice_Type(t *testing.T) {
	// Setup cache
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Error("error creating device cache:", err)
		t.Fail()
		return
	}

	// when models are nested

	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12B466 [FBDV/iPhone7,2]"

	device := NewDevice(cache, userAgent)

	if device.Type() != "smartphone" {
		t.Errorf("expected device to be of 'smartphone' type but was '%s'", device.Type())
	}

	// when models are NOT nested

	userAgent = "AIRNESS-AIR99/REV 2.2.1/Teleca Q03B1"

	device = NewDevice(cache, userAgent)

	if device.Type() != "feature phone" {
		t.Errorf("expected device to be of 'smartphone' type but was '%s'", device.Type())
	}

	// When it cannot find a device type

	userAgent = "UNKNOWN MODEL TYPE"

	device = NewDevice(cache, userAgent)

	if device.Type() != "" {
		t.Error("expected device to return empty type but got:", device.Type())
	}

	// device not specified in nested block

	userAgent = "Mozilla/5.0 (Linux; Android 4.4.2; es-us; SAMSUNG SM-G900F Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko)"

	device = NewDevice(cache, userAgent)

	if device.Type() != "smartphone" {
		t.Error("should fall back to top-level device type, but got:", device.Type())
	}

	// concrete device types

	for _, test := range typeTests {
		device := NewDevice(cache, test.UserAgent)

		if device.Name() != test.ExpectedName {
			t.Errorf("test '%s': expected Name() to return '%s' but returned '%s'.\n", test.Description, test.ExpectedName, device.Name())
		}

		if device.Brand() != test.ExpectedBrand {
			t.Errorf("test '%s': expected Brand() to return '%s' but returned '%s'.", test.Description, test.ExpectedBrand, device.Brand())
		}

		if device.Type() != test.ExpectedType {
			t.Errorf("test '%s': expected Type() to return '%s' but returned '%s'.\n", test.Description, test.ExpectedType, device.Type())
		}
	}
}
