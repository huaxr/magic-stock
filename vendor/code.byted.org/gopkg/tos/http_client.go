package tos

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const clusterSep = "$"
const connectionPoolSep = "$$" // see  https://bytedance.feishu.cn/space/doc/doccnd33ItuvZOmoViIV7Kk5Ifg#
const maxConnRetryTimes = 2
const connectionPoolReuseSpan = time.Duration(60) * time.Second

var dialer = net.Dialer{Timeout: 3 * time.Second}
var availableEndpoints = []string{
	"tos-cn-north.byted.org",
	"tos-cn-north-lf.byted.org",
	"tos-cn-north-hl.byted.org",
	"tos-cn-north-lq.byted.org",
	"tos-cn-north-1.byted.org",
	"tos-cn-north-90.byted.org",
}

// 1. use only conn succ ratio to decide quality
//    when get one addr, connect start, connect done.
// 2. how to do retry, how to count retry times?
//    use ctx

type httpClient struct {
	client  http.Client
	addrMan *addrManager
	cluster string
}

func isIPAddr(s string) bool {
	h, _, _ := net.SplitHostPort(s)
	return h != "" && net.ParseIP(h) != nil
}

func isPsm(s string) bool {
	return strings.Contains(s, defaultServiceName)
}

func isEndpointValidDomain(addr string) bool {
	match := false
	for _, endpoint := range availableEndpoints {
		if strings.HasSuffix(addr, endpoint) {
			match = true
		}
	}
	return match
}

// NewHttpClient returns HttpClient with `http://{YOUR_SERVICE}/path/to/your/api` support
func newHttpClient(cluster, idc, endpoint, serviceName string) (*httpClient, error) {
	var addrMan *addrManager = nil
	var err error = nil
	if endpoint == "" {
		addrMan, err = newAddrManager(cluster, idc, serviceName)
		if err != nil {
			return nil, err
		}
	} else if !isEndpointValidDomain(endpoint) {
		return nil, errors.New("endpoint is not a valid domain")
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if idx := strings.Index(addr, connectionPoolSep); idx > 0 {
				addr = addr[:idx]
			}
			if idx := strings.Index(addr, clusterSep); idx > 0 {
				addr = addr[:idx] // rm cluster name
			}
			if isIPAddr(addr) {
				return dialer.DialContext(ctx, network, addr)
			}
			if isEndpointValidDomain(addr) {
				return dialer.DialContext(ctx, network, addr+":80")
			}
			var err error
			var conn net.Conn
			for i := 0; i < maxConnRetryTimes; i++ {
				addr := addrMan.getAddr()
				conn, err = dialer.DialContext(ctx, network, addr)
				if err == nil {
					addrMan.cntSucc(addr)
					return conn, nil
				} else {
					if strings.Contains(err.Error(), "connection refused") {
						addrMan.fastCntFail(addr)
					} else {
						addrMan.cntFail(addr)
					}
				}
			}

			return nil, fmt.Errorf("retry %d times, still cannot connect to server, last conn err: %s", maxConnRetryTimes, err.Error())
		},
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     10 * time.Second,
		DisableCompression:  true,
	}

	return &httpClient{
		client:  http.Client{Transport: transport},
		addrMan: addrMan,
		cluster: cluster,
	}, nil
}

func (c *httpClient) do(req *http.Request) (*http.Response, error) {
	if !isIPAddr(req.URL.Host) {
		if c.cluster != "" {
			// add cluster name as a part of host
			//	in order to identity http connection pool
			req.URL.Host += clusterSep + c.cluster
		}
		if isPsm(req.URL.Host) {
			req.URL.Host += connectionPoolSep + fmt.Sprintf("%d", time.Now().Truncate(connectionPoolReuseSpan).Unix())
		}
	}

	return c.client.Do(req)
}
