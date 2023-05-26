package versionextractor

import (
	"log"
	"strings"

	"github.com/gijsbers/go-pcre"
	"github.com/robicode/device-detector/extractor"
	"github.com/robicode/device-detector/util"
	"github.com/robicode/version"
)

var (
	MajorVersion2 = version.New2("2.0")
	MajorVersion3 = version.New2("3.0")
	MajorVersion4 = version.New2("4.0")
	MajorVersion8 = version.New2("8.0")
)

type Version struct {
	Regex          string
	Version        string
	compiledRegexp pcre.Regexp
	compileError   error
	compiled       bool
}

type VersionExtractor struct {
	_regex     string
	_userAgent string
	_version   string
	_versions  []Version
}

func New(userAgent string, regex string, version string, versions []Version) *VersionExtractor {
	return &VersionExtractor{
		_regex:     regex,
		_userAgent: userAgent,
		_version:   version,
		_versions:  versions,
	}
}

func (v *VersionExtractor) Call() string {
	simpleVersion := extractor.New(v._userAgent, v._regex, v._version).Call()
	simpleVersion = strings.ReplaceAll(simpleVersion, "_", ".")
	simpleVersion = strings.TrimSuffix(simpleVersion, ".")
	if simpleVersion != "" {
		return simpleVersion
	}
	return v.osVersionByRegexes()
}

func (v *VersionExtractor) osVersionByRegexes() string {
	if v._versions == nil || len(v._versions) == 0 {
		return ""
	}

	for _, version := range v._versions {
		if !version.compiled && version.compileError == nil {
			re, err := pcre.Compile(version.Regex, pcre.CASELESS)
			if err != nil {
				version.compileError = err
				log.Println(err)
			}
			version.compiled = true
			version.compiledRegexp = re
		}

		if version.compiled {
			matcher := version.compiledRegexp.MatcherString(v._userAgent, 0)
			if matcher.Matches() {
				re := pcre.MustCompile(`\$(\d)`, pcre.CASELESS)
				s := util.EGSub(version.Version, re, func(s1 string, i int, s2 []string) string {
					return matcher.GroupString(1)
				})
				return strings.TrimSpace(strings.ReplaceAll(s, "_", "."))
			}
		}
	}
	return ""
}
