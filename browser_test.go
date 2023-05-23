package devicedetector

import "testing"

func TestBrowser_MobileOnlyBrowser(t *testing.T) {
	if !mobileOnlyBrowser("OC") {
		t.Error("expected 'OC' to be in the list of mobile only browsers")
	}

	if !mobileOnlyBrowser("AdBlock Browser") {
		t.Error("expected long name to be in mobileOnlyBrowsers")
	}
}
