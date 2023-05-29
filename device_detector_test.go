package devicedetector

import (
	"testing"
)

func TestDeviceDetector(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil || cache == nil {
		t.Errorf("error creating embesubjected cache: %v", err)
		t.Fail()
		return
	}

	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.69"

	subject := New(cache, userAgent)

	if subject.Name() != "Chrome" {
		t.Errorf("expected Name() to return 'Chrome' but got '%s'.", subject.Name())
	}

	if subject.FullVersion() != "30.0.1599.69" {
		t.Errorf("expected FullVersion() to return '30.0.1599.69' but returned '%s'", subject.FullVersion())
	}

	if subject.OSFamily() != "Mac" {
		t.Errorf("expected OSFamily() to return 'Mac' but returned '%s'", subject.OSFamily())
	}

	if subject.OSName() != "Mac" {
		t.Errorf("expected OSName() to return 'Mac' but returned '%s'", subject.OSName())
	}

	if subject.OSFullVersion() != "10.8.5" {
		t.Errorf("expected OSFullVersion() to return the correct '10.8.5' but returned '%s'", subject.OSFullVersion())
	}

	if !subject.IsKnown() {
		t.Error("expected client to be known")
	}

	if subject.IsBot() {
		t.Error("expected client not to be a bot")
	}

	if subject.BotName() != "" {
		t.Errorf("expected BotName() to return empty string but returned '%s'", subject.BotName())
	}

	// Ubuntu Linux
	userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36"

	subject = New(cache, userAgent)

	if subject.OSFamily() != "GNU/Linux" {
		t.Errorf("expected OSFamily() to return 'GNU/Linux' but returned '%s'", subject.OSFamily())
	}

	if subject.OSName() != "Ubuntu" {
		t.Errorf("expected OSName() to return 'Ubuntu' but returned '%s'", subject.OSName())
	}

	// Firefox mobile phone

	userAgent = "Mozilla/5.0 (Android 7.0; Mobile; rv:53.0) Gecko/53.0 Firefox/53.0"
	subject = New(cache, userAgent)

	if subject.DeviceType() != "smartphone" {
		t.Errorf("expected DeviceType() to return 'smartphone' but returned '%s'", subject.DeviceType())
	}

	// Firefox mobile tablet

	userAgent = "Mozilla/5.0 (Android 6.0.1; Tablet; rv:47.0) Gecko/47.0 Firefox/47.0"
	subject = New(cache, userAgent)

	if subject.DeviceType() != "tablet" {
		t.Errorf("expected DeviceType() to return 'tablet' but returned '%s'", subject.DeviceType())
	}

	// unknown user agent

	userAgent = "garbage123"
	subject = New(cache, userAgent)

	if subject.Name() != "" {
		t.Errorf("expected Name() to be blank but returned '%s'", subject.Name())
	}

	if subject.FullVersion() != "" {
		t.Errorf("expected FullVersion() to be blank but returned '%s'", subject.FullVersion())
	}

	if subject.OSName() != "" {
		t.Errorf("expected OSName() to be blank but returned '%s'", subject.OSName())
	}

	if subject.OSFullVersion() != "" {
		t.Errorf("expected OSFullVersion() to return '' but returned '%s'", subject.OSFullVersion())
	}

	if subject.IsKnown() {
		t.Error("expected IsKnown() to be false but was true")
	}

	if subject.IsBot() {
		t.Error("expected IsBot() to be false but was true")
	}

	if subject.BotName() != "" {
		t.Errorf("expected BotName() to be blank but returned '%s'", subject.BotName())
	}

	// Bot

	userAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	subject = New(cache, userAgent)

	if subject.Name() != "" {
		t.Errorf("expected Name() to return '' but returned '%s'", subject.Name())
	}

	if subject.FullVersion() != "" {
		t.Errorf("expected FullVersion() to be '' but was '%s'", subject.FullVersion())
	}

	if subject.OSName() != "" {
		t.Errorf("expected OSName() to be '' but was '%s'", subject.OSName())
	}

	if subject.OSFullVersion() != "" {
		t.Errorf("expected OSFullVersion() to be '' but was '%s'", subject.OSFullVersion())
	}

	if subject.IsKnown() {
		t.Error("expected IsKnown() to be false but was true")
	}

	if !subject.IsBot() {
		t.Error("expected IsBot() to be true but was false")
	}

	if subject.BotName() != "Googlebot" {
		t.Errorf("expected BotName() to be 'Googlebot' but was '%s'", subject.BotName())
	}
}
