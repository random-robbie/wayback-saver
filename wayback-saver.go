package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	var url string
	flag.StringVar(&url, "url", "", "Url to Save to Wayback machine")

	flag.Parse()

	if len(url) == 0 {
		fmt.Println("Usage: wayback-saver -url http://example.com")
		flag.PrintDefaults()
		os.Exit(1)
	}

	headers := map[string]string{
		"Host":            "firefox-api.archive.org",
		"Cookie":          "donation-identifier=d3f62d46b2bc196ba05db174992496eb",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:103.0) Gecko/20100101 Firefox/103.0 Wayback_Machine_Firefox/3.1",
		"Accept":          "text/html,application/xhtml+xml,application/xml",
		"Accept-Language": "en-US,en;q=0.5",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Content-Length":  "362",
		"Origin":          "moz-extension://f4e4a1e5-acd9-4e04-a0b0-252ecf93ea70",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"Te":              "trailers",
	}
	var data = []byte("capture_all=1&url=" + url + "")
	httpRequest("https://firefox-api.archive.org/save/?capture_all=1&url="+url+"", "POST", data, headers)

}

func httpRequest(targetUrl string, method string, data []byte, headers map[string]string) *http.Response {

	request, error := http.NewRequest(method, targetUrl, bytes.NewBuffer(data))
	for k, v := range headers {
		request.Header.Set(k, v)

	}

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}
	response, error := client.Do(request)
	defer response.Body.Close()

	if error != nil {
		panic(error)
	}

	body, _ := ioutil.ReadAll(response.Body)
	findsaving := "Saving page"
	bodystr := string(body)
	if strings.Contains(bodystr, findsaving) {
		fmt.Println("Page is Being Saved")
	}
	syntaxstr := "URL syntax is not valid"
	if strings.Contains(syntaxstr, findsaving) {
		fmt.Println("URL syntax is not valid")
	}
	unablestr := "Cannot resolve host"
	if strings.Contains(unablestr, findsaving) {
		fmt.Println("Cannot resolve host")
	}

	return response
}
