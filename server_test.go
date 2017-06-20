package sdproxy

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
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

	apiWeb := NewLocation("/api/web", NewUpstream(&httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8091"
	}}, &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8092"
	}}))

	server := NewServer(web, apiWeb, api)

	location := server.getLocation("/api/web/v1/calendar")
	assert.Equal(t, apiWeb, location)

	location = server.getLocation("/api/v1/calendar")
	assert.Equal(t, api, location)

	location = server.getLocation("/anything/else")
	assert.Equal(t, web, location)
}

func TestGetLocationNil(t *testing.T) {
	api := NewLocation("/api", NewUpstream(&httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8091"
	}}, &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8092"
	}}))

	server := NewServer(api)
	assert.Nil(t, server.getLocation("/anything/else"))
}

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

	apiWeb := NewLocation("/api/web", NewUpstream(&httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8091"
	}}, &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8092"
	}}))

	server := NewServer(web, apiWeb, api)
	server.sortLocations()
	sortedApiWeb := server.locations[0]
	sortedApi := server.locations[1]
	sortedWeb := server.locations[2]

	assert.Equal(t, apiWeb, sortedApiWeb)
	assert.Equal(t, web, sortedWeb)
	assert.Equal(t, api, sortedApi)
}
