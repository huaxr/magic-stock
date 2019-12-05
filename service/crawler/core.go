// @Contact:    huaxinrui
// @Time:       2019/7/5 下午5:13

package crawler

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	PROXY = "http://bj-rd-proxy.byted.org:3128"
)

func (c *Crawler) NewDocumentWithProxy(uri string) (*goquery.Document, error) {
	// Load the URL
	proxy, _ := url.Parse(PROXY)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}
	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(res.Body)
}

func (c *Crawler) NewDocument(url string) (*goquery.Document, error) {
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(res.Body)
}
