charsetx
========

[![GoDoc](https://godoc.org/github.com/philipjkim/charsetx?status.svg)](https://godoc.org/github.com/philipjkim/charsetx) [![Go Report Card](https://goreportcard.com/badge/github.com/philipjkim/charsetx)](https://goreportcard.com/report/github.com/philipjkim/charsetx) [![Build Status](https://travis-ci.org/philipjkim/charsetx.svg)](https://travis-ci.org/philipjkim/charsetx)

charsetx detects charset encoding of an HTML document, and convert a non-UTF8 page body to UTF8 string.

There are 3 steps for charset detection:

1. Return the result of [charset.DetermineEncoding()](https://godoc.org/golang.org/x/net/html/charset#DetermineEncoding) if `certain` is true.
2. Return the result of [chardet.Detector.DetectBest()](https://godoc.org/github.com/saintfish/chardet#Detector.DetectBest) if `Confidence` is 100.
3. Return charset in `Content-Type` meta tag if exists.

If all 3 steps fails, it returns error.

## Install

    go get -u github.com/philipjkim/charsetx


## Example

### Getting UTF-8 string of body for given URL

```go
// Invalid UTF-8 characters are discarded to give a result 
// rather than giving error 
// if the second bool param is set to true.
r, err := charsetx.GetUTF8BodyFromURL("http://www.godoc.org", false)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Println(r)
```

If you want to reuse or customize `*http.Client` instead of `http.DefaultClient`, 
use [GetUTF8Body()](https://godoc.org/github.com/philipjkim/charsetx#GetUTF8Body).


### Getting the charset of given URL

```go
client := http.DefaultClient
resp, err := client.Get("http://www.godoc.org")
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
```


## Related Projects

* [golang.org/x/net/html/charset](https://godoc.org/golang.org/x/net/html/charset)
* [github.com/saintfish/chardet](https://godoc.org/github.com/saintfish/chardet)
* [github.com/djimenez/iconv-go](https://godoc.org/github.com/djimenez/iconv-go)


## License

[MIT](LICENSE)
