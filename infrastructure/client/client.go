package client

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/task-done/infrastructure/log"
)

var httpClient *http.Client

func InitHttpClient() {
	var defaultTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient = &http.Client{Transport: defaultTransport}
}

func GET(url string, requestBody io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, requestBody)
	if err != nil {
		log.Error("new http GET request error:%s", err)
		return nil, err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		log.Error("do the http GET request error:%s", err)
		return nil, err
	}
	return response, nil
}
