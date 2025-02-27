package httpClient

import (
	"cw/logger"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

type Option func(*HttpClient)

func WithHttp2() Option {
	return func(hc *HttpClient) {
		transport, ok := hc.Client.Transport.(*http.Transport)
		if !ok {
			logger.GlobalLogger.Warnf("Failed create transport in options constructor")
			return
		}
		if err := http2.ConfigureTransport(transport); err != nil {
			logger.GlobalLogger.Error(err)
			return
		}
	}
}

func WithProxy(proxy string) Option {
	return func(hc *HttpClient) {
		if proxy != "" {
			proxyURL, err := url.Parse(proxy)
			if err != nil {
				logger.GlobalLogger.Errorf("failed o create proxy url: %v", err)
				return
			}
			hc.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
		}
	}
}
