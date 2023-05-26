package versionextractor

import (
	"strings"
	"testing"
)

func TestVersionExtractor_Call(t *testing.T) {
	userAgent := "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0; Avant Browser; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0)"
	regex := `Avant Browser`
	version := ""

	ve := New(userAgent, regex, version, nil)

	if ve.Call() != "" {
		t.Errorf("expected extractor without version to return empty string but got '%s'", ve.Call())
	}

	userAgent = "Mozilla/5.0 (X11; U; Linux i686; nl; rv:1.8.1b2) Gecko/20060821 BonEcho/2.0b2 (Debian-1.99+2.0b2+dfsg-1)"
	version = "BonEcho (2.0)"
	regex = `(BonEcho|GranParadiso|Lorentz|Minefield|Namoroka|Shiretoko)/(\d+[\.\d]+)`
	metaVersion := "$1 ($2)"

	result := New(userAgent, regex, metaVersion, nil).Call()

	if result != version {
		t.Errorf("regex with dynamic matching failed with result '%s' (expected '%s')", result, version)
	}

	version = version + "  "
	result = New(userAgent, regex, metaVersion, nil).Call()

	if result != strings.TrimSpace(version) {
		t.Errorf("expected Call() to strip whitespace; '%s' did not equal result '%s'", version, result)
	}

	userAgent = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)"
	regex = `MSIE.*Trident/4.0`
	version = "8.0"

	result = New(userAgent, regex, version, nil).Call()

	if result != version {
		t.Errorf("extractor with fixed version did not return proper version of '%s' but returned '%s'", version, result)
	}

	userAgent = "garbage"

	result = New(userAgent, `any`, "", nil).Call()

	if result != "" {
		t.Errorf("expected unknown user agent to return empty string but got '%s'", result)
	}
}
