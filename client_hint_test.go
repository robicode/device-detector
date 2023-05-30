package devicedetector

import (
	"net/http"
	"testing"
)

func TestClientHint(t *testing.T) {
	headers := http.Header{}
	headers.Add("sec-ch-ua", `"Opera";v="83", " Not;A Brand";v="99", "Chromium";v="98"`)
	headers.Add("sec-ch-ua-mobile", `?0`)
	headers.Add("sec-ch-ua-platform", "Windows")
	headers.Add("sec-ch-ua-platform-version", "14.0.0")

	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Error("error initialising cache:", err)
	}

	ch := NewClientHint(cache, headers)

	if ch.Platform() != "Windows" {
		t.Errorf("expected Platform() to be 'Windows' but was '%s'", ch.Platform())
	}

	if ch.PlatformVersion() != "14.0.0" {
		t.Errorf("expected PlatformVersion() to be '14.0.0' but was '%s'", ch.PlatformVersion())
	}

	if ch.BrowserName() != "Opera" {
		t.Errorf("expected BrowserName() to return 'Opera' but returned '%s'", ch.BrowserName())
	}

	testBrands := map[string]string{
		"Opera":        "83",
		" Not;A Brand": "99",
		"Chromium":     "98",
	}

	for brand, version := range ch.Brands() {
		if item, ok := testBrands[brand]; !ok {
			t.Errorf("expected '%s' to be in the list of brands", item)
		} else {
			if version != item {
				t.Errorf("expected version to be '%s' but was '%s'", item, version)
			}
		}
	}
}
