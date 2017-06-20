package sdproxy

import (
	"net/http"
	"net/http/httputil"
)

func NewReverseProxy(host string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = host
		},
	}
}
