package sdproxy

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	web := NewLocation("/", NewUpstream("127.0.0.1:8090", "127.0.0.1:8094"))
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	apiWeb := NewLocation("/api/web", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))

	server := NewServer("127.0.0.1:8080", web, apiWeb, api)
	req := &http.Request{
		URL: &url.URL{
			Path: "/api/web/v1/calendar",
		},
	}
	location := server.getLocation(req)
	assert.Equal(t, apiWeb, location)
	req = &http.Request{
		URL: &url.URL{
			Path: "/api/v1/calendar",
		},
	}
	location = server.getLocation(req)
	assert.Equal(t, api, location)
	req = &http.Request{
		URL: &url.URL{
			Path: "/anything/else",
		},
	}
	location = server.getLocation(req)
	assert.Equal(t, web, location)
}

func TestGetLocationHeader(t *testing.T) {
	userAgent := &MatchHeader{
		Name:    "User-Agent",
		Pattern: "Googlebot",
	}
	web := NewLocation("/webapi", NewUpstream("127.0.0.1:8090", "127.0.0.1:8094"))
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"), userAgent)

	server := NewServer("127.0.0.1:8080", web, api)
	// match path but not header
	req := &http.Request{
		URL: &url.URL{
			Path: "/api/web/v1/calendar",
		},
		Header: http.Header{"User-Agent": []string{"Googlebo"}},
	}
	assert.Nil(t, server.getLocation(req))
	req = &http.Request{
		URL: &url.URL{
			Path: "/api/web/v1/calendar",
		},
		Header: http.Header{"User-Agent": []string{"Googlebot"}},
	}
	location := server.getLocation(req)
	assert.Equal(t, api, location)
}

func TestGetLocationNil(t *testing.T) {
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	server := NewServer("127.0.0.1:8080", api)
	req := &http.Request{
		URL: &url.URL{
			Path: "/anything/else",
		},
	}
	assert.Nil(t, server.getLocation(req))
}

func TestLocationByPath(t *testing.T) {
	web := NewLocation("/", NewUpstream("127.0.0.1:8090", "127.0.0.1:8094"))
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	apiWeb := NewLocation("/api/web", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	apiWebHeader := NewLocation("/api/web", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"), &MatchHeader{})

	server := NewServer("127.0.0.1:8080", web, apiWeb, api, apiWebHeader)
	server.sortLocations()
	sortedWebHeader := server.locations[0]
	sortedApiWeb := server.locations[1]
	sortedApi := server.locations[2]
	sortedWeb := server.locations[3]
	assert.Equal(t, apiWebHeader, sortedWebHeader)
	assert.Equal(t, apiWeb, sortedApiWeb)
	assert.Equal(t, web, sortedWeb)
	assert.Equal(t, api, sortedApi)
}

func TestServerAddr(t *testing.T) {
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	server := NewServer("127.0.0.1:8080", api)
	assert.Equal(t, "127.0.0.1:8080", server.addr)
	server.SetAddr("127.0.0.1:9090")
	assert.Equal(t, "127.0.0.1:9090", server.addr)
}
