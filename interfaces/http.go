package interfaces

import (
	"crypto/tls"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// newInsecureHttpClient creates an HTTP client with TLS verification disabled.
func NewInsecureHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}
