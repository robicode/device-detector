package devicedetector

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type TestBot struct {
	Name     string `yaml:"name"`
	Category string `yaml:"category"`
	Url      string `yaml:"url"`
}

type BotFixture struct {
	UserAgent string  `yaml:"user_agent"`
	Bot       TestBot `yaml:"bot"`
}

var botFixtures []BotFixture

func TestBotWithFixtures(t *testing.T) {
	cache, err := NewEmbeddedCache()
	if err != nil {
		t.Error("error creating cache:", err)
		t.Fail()
		return
	}

	if cache == nil {
		t.Error("cache created but is nil?")
		t.Fail()
		return
	}

	data, err := os.ReadFile("fixtures/bots.yml")
	if err != nil {
		t.Error("error loading bot fixtures:", err)
		t.Fail()
		return
	}

	err = yaml.Unmarshal(data, &botFixtures)
	if err != nil {
		t.Error("error parsing fixtures:", err)
		t.Fail()
		return
	}

	for _, fixture := range botFixtures {
		bot := NewBot(cache, fixture.UserAgent)
		if !bot.IsBot() {
			t.Errorf("expected bot '%s' to be a bot", bot.Name())
		}

		if bot.Name() != fixture.Bot.Name {
			t.Errorf("expected bot name to be '%s' but was '%s'", fixture.Bot.Name, bot.Name())
		}
	}
}
