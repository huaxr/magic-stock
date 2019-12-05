// @Contact:    huaxinrui
// @Time:       2019/7/5 下午5:13

package core

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 都是历史记录的静态数据
func NewDocumentWithProxy(uri string) (*goquery.Document, error) {
	// Load the URL
	proxy, _ := url.Parse("http://bj-rd-proxy.byted.org:3128")
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}
	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func NewDocument(url string) (*goquery.Document, error) {
	// Load the URL
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	return goquery.NewDocumentFromReader(res.Body)
}
