package util

import (
	"testing"

	"github.com/gijsbers/go-pcre"
)

func TestUtil_EGSub(t *testing.T) {
	re, err := pcre.Compile(`\$(\d)`, 0)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if s := EGSub("$1", re, func(matched string, match int) string {
		return "Test"
	}); s != "Test" {
		t.Errorf("[utils]: expected EGSub to return 'Test' but got '%s'", s)
		t.Fail()
		return
	}

	re = pcre.MustCompile(`[aeiou]`, 0)

	if s := EGSub("hello", re, func(matched string, match int) string {
		return "*"
	}); s != "h*ll*" {
		t.Errorf("[utils]: s != 'h*ll*'; s == '%s'", s)
		t.Fail()
		return
	}
}

func TestUtil_FixupRegex(t *testing.T) {
	if FixupRegex(`XYZZY`) != `(?:^|[^A-Z0-9\-_]|[^A-Z0-9\-]_|sprd-|MZ-)(?:XYZZY)` {
		t.Errorf("expected BuildRegex() to include the passed in regex plus the base regex but got `%s`", FixupRegex(`XYZZY`))
		t.Fail()
		return
	}
}
