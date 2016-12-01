package csdetect

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	URL         string
	ExpectedStr string
}

func TestGetUTF8BodyForUT8Pages(t *testing.T) {
	data := []TestData{
		TestData{URL: "http://www.godoc.org", ExpectedStr: "GoDoc"},
		TestData{URL: "http://www.kakaocorp.com", ExpectedStr: "카카오"},
	}
	for _, d := range data {
		r, err := GetUTF8Body(d.URL)
		assert.Nil(t, err)
		assert.Contains(t, r, d.ExpectedStr)
	}
}

func TestGetUTF8BodyForNonUT8KoreanPages(t *testing.T) {
	data := []TestData{
		TestData{URL: "http://blog.naver.com/tt820613/220017990859", ExpectedStr: "파라다이스"},
		TestData{URL: "http://blog.naver.com/yeseul961/220649621993", ExpectedStr: "블로그"},
		TestData{URL: "http://piggohome.com", ExpectedStr: "집으로돼지"},
	}
	for _, d := range data {
		r, err := GetUTF8Body(d.URL)
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
