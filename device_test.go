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
