package sdproxy

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationByPath(t *testing.T) {
	web := NewLocation("/", NewUpstream(&httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8090"
	}}, &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8094"
	}}))
	api := NewLocation("/api", NewUpstream(&httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8091"
	}}, &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8092"
	}}))

	server := NewServer(web, api)
	server.sortLocations()
	sortedWeb := server.locations[1]
	sortedApi := server.locations[0]
	assert.Equal(t, web, sortedWeb)
	assert.Equal(t, api, sortedApi)
}
