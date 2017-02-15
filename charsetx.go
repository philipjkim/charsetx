// Package charsetx provides functions for detecting charset encoding of
// an HTML document and UTF-8 conversion.
package charsetx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// GetUTF8Body returns string of body.
// If charset for body is detected as non-UTF8,
// this function converts to UTF-8 string and return it.
// Illegal byte sequences are silently discarded
// if ignoreInvalidUTF8Chars is set to true.
func GetUTF8Body(body []byte, contentType string,
	ignoreInvalidUTF8Chars bool) (string, error) {
	// Detect charset.
	cs, err := DetectCharset(body, contentType)
	if err != nil {
		if !ignoreInvalidUTF8Chars {
			return "", err
		}
		cs = "utf-8"
	}

	// Remove utf8.RuneError if ignoreInvalidUTF8Chars is true.
	bs := string(body)
	if ignoreInvalidUTF8Chars {
		if !utf8.ValidString(bs) {
			v := make([]rune, 0, len(bs))
			for i, r := range bs {
				if r == utf8.RuneError {
					_, size := utf8.DecodeRuneInString(bs[i:])
					if size == 1 {
						continue
					}
				}
				v = append(v, r)
			}
			bs = string(v)
		}
	}

	// Convert body.
	converted, err := iconv.ConvertString(bs, cs, "utf-8")
	if err != nil && !strings.Contains(converted, "</head>") {
		return "", err
	}

	return converted, nil
}

// GetUTF8BodyFromURL returns response body of urlStr as string.
// It converts charset to UTF-8 if original charset is non UTF-8.
// Illegal byte sequences are silently discarded if ignoreIBS is set to true.
func GetUTF8BodyFromURL(urlStr string, ignoreIBS bool) (string, error) {
	client := http.DefaultClient
	resp, err := client.Get(urlStr)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Detect charset.
	return GetUTF8Body(byt, resp.Header.Get("Content-Type"), ignoreIBS)
}

/*
DetectCharset returns charset for contents of urlStr by 3 steps:

   1. Use the result of charset.DetermineEncoding() if certain is true.
   2. Use the result of chardet.Detector.DetectBest() if Confidence is 100.
   3. Use charset in `Content-Type` meta tag if exists.

If no charset is detected by 3 steps, it returns error.

Argument body is response body as byte[],
and contentType is `Content-Type` response header value.
*/
func DetectCharset(body []byte, contentType string) (string, error) {
	// 1. Use charset.DetermineEncoding
	_, name, certain := charset.DetermineEncoding(body, contentType)
	if certain {
		return name, nil
	}

	// Handle uncertain cases
	// 2. Use chardet.Detector.DetectBest
	r, err := chardet.NewHtmlDetector().DetectBest(body)
	if err != nil {
		return "", err
	}
	if r.Confidence == 100 {
		return strings.ToLower(r.Charset), nil
	}

	// 3. Parse meta tag for Content-Type
	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	doc := goquery.NewDocumentFromNode(root)
	var csFromMeta string
	doc.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// <meta http-equiv="Content-Type" content="text/html; charset=MS949"/>
		if c, exists := s.Attr("content"); exists && strings.Contains(c, "charset") {
			if _, params, err := mime.ParseMediaType(c); err == nil {
				if cs, ok := params["charset"]; ok {
					csFromMeta = strings.ToLower(cs)
					// Handle Korean charsets.
					if csFromMeta == "ms949" || csFromMeta == "cp949" {
						csFromMeta = "euc-kr"
					}
					return false
				}
			}
			return true
		}

		// <meta charset="utf-8"/>
		if c, exists := s.Attr("charset"); exists {
			csFromMeta = c
			return false
		}

		return true
	})

	if csFromMeta == "" {
		return "", fmt.Errorf("failed to detect charset")
	}

	return csFromMeta, nil
}
