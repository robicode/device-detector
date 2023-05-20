package modelextractor

import "testing"

func TestModelExtractor_Call(t *testing.T) {
	// When no dynamic match is found
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4"
	regex := `(?:Apple-)?iPhone ?(3GS?|4S?|5[CS]?|6(:? Plus)?)?`

	extractor := New(userAgent, regex, "iPhone $1")

	if extractor.Call() != "iPhone" {
		t.Errorf("expected 'iPhone' but got '%s'", extractor.Call())
	}

	// when a dynamic match is found
	userAgent = "Mozilla/5.0 (iPhone 5S; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4"
	extractor = New(userAgent, regex, "iPhone $1")

	if extractor.Call() != "iPhone 5S" {
		t.Errorf("expected 'iPhone 5S' but got '%s'", extractor.Call())
	}

	// when matching against static model

	userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12A365 Weibo (iPhone7,2)"
	regex = `(?:Apple-)?iPhone7[C,]2`

	extractor = New(userAgent, regex, "iPhone 6")

	if extractor.Call() != "iPhone 6" {
		t.Errorf("expected 'iPhone 6' but got '%s'", extractor.Call())
	}
}
