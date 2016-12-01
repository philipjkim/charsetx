charsetx
========

[![GoDoc](https://godoc.org/github.com/philipjkim/charsetx?status.svg)](https://godoc.org/github.com/philipjkim/charsetx) [![Go Report Card](https://goreportcard.com/badge/github.com/philipjkim/charsetx)](https://goreportcard.com/report/github.com/philipjkim/charsetx) [![Build Status](https://travis-ci.org/philipjkim/charsetx.svg)](https://travis-ci.org/philipjkim/charsetx)

charsetx detects charset encoding of an HTML document, and convert a non-UTF8 page body to UTF8 string.

There are 3 steps for charset detection:

1. Return the result of [charset.DetermineEncoding()](https://godoc.org/golang.org/x/net/html/charset#DetermineEncoding) if `certain` is true.
2. Return the result of [chardet.Detector.DetectBest()](https://godoc.org/github.com/saintfish/chardet#Detector.DetectBest) if `Confidence` is 100.
3. Return charset in `Content-Type` meta tag if exists.
4. Else, return error.

Install
-------

    go get -u github.com/philipjkim/charsetx


Example
-------

Getting UTF-8 string of body for given URL:

```go
package main

import (
    "fmt"
    "github.com/philipjkim/charsetx"
)

func main() {
    r, err := charsetx.GetUTF8Body("http://www.godoc.org")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(r)
}
```

Getting the charset of given URL:

```go
package main

import (
	"io/ioutil"
	"net/http"
    "fmt"
    "github.com/philipjkim/charsetx"
)

func main() {
	client := http.DefaultClient
	resp, err := client.Get(u)
	if err != nil {
		fmt.Println(err)
        return
	}

	defer resp.Body.Close()
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
        return
	}

    cs, err := charsetx.DetectCharset(byt, resp.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println(err)
        return
	}

    fmt.Println(cs)
}
```

Related Projects
----------------

* [golang.org/x/net/html/charset](https://godoc.org/golang.org/x/net/html/charset)
* [github.com/saintfish/chardet](https://godoc.org/github.com/saintfish/chardet)
* [github.com/qiniu/iconv](https://godoc.org/github.com/qiniu/iconv)


License
-------

[MIT](LICENSE)