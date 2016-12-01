package charsetx

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	URL         string
	ExpectedStr string
}

func TestGetUTF8BodyForUT8Pages(t *testing.T) {
	client := http.DefaultClient
	data := []TestData{
		{URL: "http://www.godoc.org", ExpectedStr: "GoDoc"},
		{URL: "http://www.kakaocorp.com", ExpectedStr: "카카오"},
	}
	for _, d := range data {
		r, err := GetUTF8Body(d.URL, client)
		assert.Nil(t, err)
		assert.Contains(t, r, d.ExpectedStr)
	}
}

func TestGetUTF8BodyForNonUT8KoreanPages(t *testing.T) {
	client := getCustomHTTPClient()
	data := []TestData{
		{URL: "http://blog.naver.com/tt820613/220017990859", ExpectedStr: "파라다이스"},
		{URL: "http://blog.naver.com/yeseul961/220649621993", ExpectedStr: "블로그"},
		{URL: "http://piggohome.com", ExpectedStr: "집으로돼지"},
		{URL: "https://open.sookmyung.ac.kr/ht_ml/w_02ed/2100.php", ExpectedStr: "숙명"},
		{URL: "https://httpd.apache.org/docs/current/upgrading.html", ExpectedStr: "Apache"},
		{URL: "http://www.skhappiness.org/webzine/vol07/sub/direction.jsp", ExpectedStr: "행복나눔재단"},
		{URL: "https://moti.or.kr/jsp/release/realname/List.jsp", ExpectedStr: "국방"},
		{URL: "http://m.newstown.co.kr/news/articleView.html?idxno=269864", ExpectedStr: "대선"},
		//{URL: "", ExpectedStr: ""},
	}
	for _, d := range data {
		r, err := GetUTF8Body(d.URL, client)
		assert.Nil(t, err)
		assert.Contains(t, r, d.ExpectedStr)
	}
}

func TestGetUTF8BodyWithDefaultClientForUT8Pages(t *testing.T) {
	data := []TestData{
		{URL: "http://www.godoc.org", ExpectedStr: "GoDoc"},
		{URL: "http://www.kakaocorp.com", ExpectedStr: "카카오"},
	}
	for _, d := range data {
		r, err := GetUTF8BodyWithDefaultClient(d.URL)
		assert.Nil(t, err)
		assert.Contains(t, r, d.ExpectedStr)
	}
}

func TestDetectCharsetForUTF8Pages(t *testing.T) {
	urls := []string{"http://www.godoc.org", "http://www.kakaocorp.com"}
	for _, u := range urls {
		b, c, err := getBodyAndContentType(u)
		assert.Nil(t, err)

		cs, err := DetectCharset(b, c)
		assert.Nil(t, err)
		assert.Equal(t, "utf-8", cs)
	}
}

func TestDetectCharsetForNonUTF8KoreanPages(t *testing.T) {
	urls := []string{
		"http://blog.naver.com/tt820613/220017990859",
		"http://blog.naver.com/yeseul961/220649621993",
		"http://piggohome.com",
	}
	for _, u := range urls {
		b, c, err := getBodyAndContentType(u)
		assert.Nil(t, err)

		cs, err := DetectCharset(b, c)
		assert.Nil(t, err)
		assert.Equal(t, "euc-kr", cs)
	}
}

func getBodyAndContentType(u string) ([]byte, string, error) {
	client := http.DefaultClient
	resp, err := client.Get(u)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// detect charset
	return byt, resp.Header.Get("Content-Type"), nil
}

func getCustomHTTPClient() *http.Client {
	timeout := 8 * time.Second
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	httpClient := &http.Client{Timeout: timeout, Transport: tr}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		if len(via) == 0 {
			return nil
		}
		for attr, val := range via[0].Header {
			if _, ok := req.Header[attr]; !ok {
				req.Header[attr] = val
			}
		}
		return nil
	}

	return httpClient
}
