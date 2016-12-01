package charsetx

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"bytes"

	"mime"

	"github.com/PuerkitoBio/goquery"
	"github.com/qiniu/iconv"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// GetUTF8Body returns response body of urlStr as string.
// It converts charset to UTF-8 if original charset is non UTF-8.
func GetUTF8Body(urlStr string) (string, error) {
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
	ct := resp.Header.Get("Content-Type")
	cs, err := DetectCharset(byt, ct)
	if err != nil {
		return "", err
	}

	// Convert body.
	cd, err := iconv.Open("utf-8", cs)
	if err != nil {
		return "", err
	}
	defer cd.Close()

	converted := cd.ConvString(string(byt))

	// TODO: Verify if broken chars exists.

	return converted, nil
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
		c, exists := s.Attr("content")
		if exists && strings.Contains(c, "charset") {
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
			fmt.Println(c)
		}
		return true
	})

	if csFromMeta == "" {
		return "", fmt.Errorf("Failed to detect charset")
	}

	return csFromMeta, nil
}
