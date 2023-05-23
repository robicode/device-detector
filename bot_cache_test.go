package devicedetector

import "testing"

func Test_NewEmbeddedBotCache(t *testing.T) {
	cache, err := NewEmbeddedBotCache()
	if err != nil {
		t.Error("error creating embedded bot cache:", err)
		t.Fail()
		return
	}

	if len(cache.bots) == 0 {
		t.Error("expected to have some bots in the cache")
	}
}
