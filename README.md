# DeviceDetector: Robust Device and Browser Detection

`DeviceDetector` is a precise and fast user agent parser and device detector written in Go, backed by the largest and most up-to-date user agent database. It is a port of the [device_detector](https://rubygems.org/gems/device_detector) Rubygem, which is in turn a port of the [Universal Device Detection Library](https://github.com/piwik/device-detector).

`DeviceDetector` will parse any user agent and detect the browser, operating system, device used (desktop, tablet, mobile, tv, cars, console, etc.), brand and model. It detects thousands of user agent strings, even from rare and obscure browsers and devices.

`DeviceDetector` provides two caching modules: `EmbeddedCache` provides a cache based on the data set embedded in the package, while `FilesystemCache` allows you to use a data set shipped alongside your application.

## A Note About the Code and Regular Expressions

This port does not aspire to be a one-to-one copy from the original code nor the Ruby gem, but rather an adaptation for the Go language.

Still, our goal is to use the original, unchanged regex yaml files, in order to mutually benefit from updates and pull request to both the original and the ported versions. This is probably not the cleanest code, and could perhaps be implemented better, and I welcome all pull requests.

Because Go uses the `re2` regular expression parser, I was not able to get most of the regular expressions working with the standard library. For example, regular expressions in the `device` subtree had about a ~480 success rate out of 1200+ regexes. As such, I opted to use [this pcre library wrapper](https://github.com/gijsbers/go-pcre) as a replacement.

I know some would rather not add an extra dependency to their code, instead sticking to standard library features, but I am not a regular expressions guru, and thus do not even know if it is possible to programmatically modify the regexps to conform to `re2`. 

I am also aware of the timing ramifications of using a replacement parser (Go provides certain timing guarantees with respect to `re2`), and I don't know if it really matters in a non-security context such as user agent parsing. If your thoughts differ, perhaps moving this functionality to a background process would be appropriate. In any case, I am open to discussions on how to solve this using `re2`; just submit an issue.

## Installation

```bash
go get github.com/robicode/device-detector
```

```go
import (
  ...
  "github.com/robicode/device-detector"
)
```

## Usage

First, create a cache that can be used for any number of device detections:

```go
// For embedded data cache: No need to distribute YAML files
// with your application.
cache, err := devicedetector.NewEmbeddedCache()
if err != nil {
  ...
}

// For filesystem data cache: You'll need to distribute the `regexes` directory
// with your application and point to it when creating the cache:
cache, err := devicedetector.NewFilesystemCache(filepath.Join(getApplicationPath(), "regexes"))
if err != nil {
  ...
}
```

In either case, `NewXXCache` will return a `*Cache` that you can pass to `New()` to avoid reloading and recompiling the regexes each
time you want to check a new UserAgent.

The cache will compile the regexps the first time they are needed.

Pass the cache to any `DeviceDetector` instances you create:

```go
userAgent := "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.17 Safari/537.36"

client, err := devicedetector.New(cache, userAgent)
if err != nil {
  ...
}

client.Name() // => "Chrome"
client.FullVersion() // => "30.0.1599.69"

client.OSName() // => "Windows"
client.OSVersion() // => "8"

// For many devices, you can also query the device name (usually the model name)
client.DeviceName() // => 'iPhone 5'
// Device types can be one of the following: desktop, smartphone, tablet,
// feature phone, console, tv, car browser, smart display, camera,
// portable media player, phablet, smart speaker, wearable, peripheral
client.DeviceType() // => 'smartphone'
```

`DeviceDetector` will return an empty string on all attributes, if the `userAgent` is unknown.
You can make a check to ensure the client has been detected:

```go
client.IsKnown() // => will return false if userAgent is unknown
```

Optionally `DeviceDetector` is using the content of `Sec-CH-UA` stored in the headers to improve the accuracy of the detection:

```go
userAgent = "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.17 Safari/537.36"
headers := http.Header{}
headers.Add("Sec-CH-UA", `"Chromium";v="106", "Brave";v="106", "Not;A=Brand";v="99"`)

client, err := devicedetector.New(cache, userAgent, headers)
// error checking snipped

client.Name() // => "Brave"
```

Same goes with `http-x-requested-with`/`x-requested-with`:

```go
userAgent := "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.17 Safari/537.36"
headers.Add("http-x-requested-with", "org.mozilla.focus")
client = devicedetector.New(cache, userAgent, headers)

client.Name() // => "Firefox Focus"
```

## License

The original code is licensed under LGPL 3, and as we use regexes from that code, this project uses the same license.

The `pcre` library we use is on GitHub with a BSD license. However, as far as I know the [code on which the fork is based](http://git.enyo.de/fw/debian/golang-pkg-pcre.git) does not appear to have a license at all. As there appear to be a few forks of this code floating around and several packages on the Go package index that use these forks, I am unsure as to whether this presents a legal issue or not, especially since it is primarily a wrapper around a C(++) API and therefore likely highly generic.

YMMV.

## Contributing

Feel free to just submit a PR or issue and I'll see to it when I can. If your feedback includes code, please be sure to include tests.