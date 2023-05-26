package util

import (
	"fmt"
	"strings"

	"github.com/gijsbers/go-pcre"
)

// only match if useragent begins with given regex or there is no letter before it
// See https://github.com/matomo-org/device-detector/blob/e7f44580a587346d74348d85322f5787d0f70363/Parser/AbstractParser.php#L304
func FixupRegex(regex string) string {
	return fmt.Sprintf(`(?:^|[^A-Z0-9\-_]|[^A-Z0-9\-]_|sprd-|MZ-)(?:%s)`, strings.ReplaceAll(regex, "/", `\/`))
}

// inStrArray checks to see if s is in the given []string array.
func InStrArray(s string, arr []string) bool {
	for _, str := range arr {
		if str == s {
			return true
		}
	}
	return false
}

// EGSub - Enumerated GSub
// EGSub is like the enumerated form of Ruby's String#gsub. It takes as parameters
// the original string, the regexp to match against, and a function to be called for
// each time the pattern is matched.
func EGSub(orig string, matcher pcre.Regexp, replacer func(string, int, []string) string) string {
	var i int

	if replacer == nil {
		return orig
	}

	for {
		m := matcher.MatcherString(orig, 0)
		if !m.Matches() {
			break
		}

		var matches []string
		if m.Groups() > 0 {
			var j int = 0

			for {
				if j > m.Groups() {
					break
				}

				matches = append(matches, m.GroupString(j))
				j++
			}
		}

		idx := m.Index()
		replacement := orig[idx[0]:idx[1]]
		orig = strings.Replace(orig, replacement, replacer(replacement, i, matches), 1)
		i++
	}

	return orig
}

func ReverseArray[T any](arr []T) []T {
	result := make([]T, len(arr))
	length := len(arr)
	for index := range arr {
		result[(length-1)-index] = arr[index]
	}
	return result
}
