package devicedetector

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v3"
)

type OSFixture struct {
	Name      string `yaml:"name"`
	ShortName string `yaml:"short_name"`
	Version   string `yaml:"version"`
	Platform  string `yaml:"platform"`
	Family    string `yaml:"family"`
}

type FOS struct {
	UserAgent string    `yaml:"user_agent"`
	OS        OSFixture `yaml:"os"`
}

func TestOSsWithFixtures(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Error("error creating embedded cache:", err)
		t.Fail()
		return
	}

	if cache == nil {
		t.Error("cache was created but is nil?")
		t.Fail()
		return
	}

	if cache.OS == nil {
		t.Error("os cache is nil")
		t.Fail()
		return
	}

	var fixtures []FOS

	data, err := ioutil.ReadFile("fixtures/oss.yml")
	if err != nil {
		t.Error("error loading fixtures:", err)
		t.Fail()
		return
	}

	err = yaml.Unmarshal(data, &fixtures)
	if err != nil {
		t.Error("error parsing fixtures:", err)
		t.Fail()
		return
	}

	for _, fixture := range fixtures {
		os := NewOS(cache, fixture.UserAgent)

		if os.Name() != fixture.OS.Name {
			t.Errorf("expected Name() to return '%s' but returned '%s'", fixture.OS.Name, os.Name())
		}

		if os.Family() != fixture.OS.Family {
			t.Errorf("os: '%s': expected Family() to return '%s' but returned '%s'", os.Name(), fixture.OS.Family, os.Family())
		}

		if os.ShortName() != fixture.OS.ShortName {
			t.Errorf("os: '%s': expected ShortName() to return '%s' but returned '%s'.", fixture.OS.Name, fixture.OS.ShortName, os.ShortName())
		}

		if os.FullVersion() != fixture.OS.Version {
			t.Errorf("os: '%s': expected FullVersion to return '%s' but returned '%s'.", os.Name(), fixture.OS.Version, os.FullVersion())
		}
	}
}
