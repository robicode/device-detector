package devicedetector

import (
	"testing"
)

func Test_NewEmbeddedDeviceCache(t *testing.T) {
	cache, err := NewEmbeddedDeviceCache()
	if err != nil {
		t.Error("error creating embedded device cache:", err)
		t.Fail()
		return
	}

	if len(cache.devices) == 0 {
		t.Error("expected to have some devices in the cache")
	}
}
