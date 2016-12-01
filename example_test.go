package charsetx_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/philipjkim/charsetx"
)

func ExampleGetUTF8Body() {
	client := http.DefaultClient
	r, err := charsetx.GetUTF8Body("http://www.godoc.org", client)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(r)
}

func ExampleGetUTF8BodyWithDefaultClient() {
	r, err := charsetx.GetUTF8BodyWithDefaultClient("http://www.godoc.org")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(r)
}

func ExampleDetectCharset() {
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
}
