package common

import (
	"io"
	"net/http"
	"time"

	"cloud-disk/internal/log"
)

var Client *HttpClient

type HttpClient struct {
	client *http.Client
}

func InitHttpClient() error {
	var defaultTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		//DialContext: (&net.Dialer{
		//        Timeout:   3 * time.Second,
		//        KeepAlive: 30 * time.Second,
		//        DualStack: true,
		//    }).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	Client = &HttpClient{
		client: &http.Client{Transport: defaultTransport},
	}
	return nil
}

func (h *HttpClient) GET(url string, requestBody io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, requestBody)
	if err != nil {
		log.Error("new http GET request error:%s", err)
		return nil, err
	}

	response, err := h.client.Do(request)
	if err != nil {
		log.Error("do the http GET request error:%s", err)
		return nil, err
	}
	return response, nil
}
