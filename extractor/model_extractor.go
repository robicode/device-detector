package extractor

import (
	"strings"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/util"
)

type Extractor struct {
	_regex     string
	_name      string
	_userAgent string
}

func New(userAgent, regex, name string) *Extractor {
	return &Extractor{
		_userAgent: userAgent,
		_name:      name,
		_regex:     regex,
	}
}

func (m *Extractor) Call() string {
	re, err := pcre.Compile(m._regex, pcre.CASELESS)
	if err != nil {
		return ""
	}

	matcher := re.MatcherString(m._userAgent, 0)
	if !matcher.Matches() {
		return ""
	}

	if matcher.Groups() == 0 {
		return m._name
	}

	re2 := pcre.MustCompile(`\$(\d)`, 0)
	s := util.EGSub(m._name, re2, func(s string, i int) string {
		return matcher.GroupString(1)
	})
	return strings.TrimSpace(s)
}
