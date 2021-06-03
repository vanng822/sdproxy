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
	locations := make([]*Location, 0)
	locations = append(locations, web)
	locations = append(locations, api)
	locations = append(locations, apiWeb)

	server := NewServer("127.0.0.1:8080", &Host{"testing.com", locations})
	req := &http.Request{
		Host: "testing.com",
		URL: &url.URL{
			Path: "/api/web/v1/calendar",
		},
	}
	location := server.getLocation(req)
	assert.Equal(t, apiWeb, location)
	req = &http.Request{
		Host: "testing.com",
		URL: &url.URL{
			Path: "/api/v1/calendar",
		},
	}
	location = server.getLocation(req)
	assert.Equal(t, api, location)
	req = &http.Request{
		Host: "testing.com",
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
	locations := make([]*Location, 0)
	locations = append(locations, NewLocation("/webapi", NewUpstream("127.0.0.1:8090", "127.0.0.1:8094")))
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"), userAgent)
	locations = append(locations, api)

	server := NewServer("127.0.0.1:8080", &Host{"testing.com", locations})
	// match path but not header
	req := &http.Request{
		Host: "testing.com",
		URL: &url.URL{
			Path: "/api/web/v1/calendar",
		},
		Header: http.Header{"User-Agent": []string{"Googlebo"}, "Host": []string{"testing.com"}},
	}
	assert.Nil(t, server.getLocation(req))
	req = &http.Request{
		Host: "testing.com",
		URL: &url.URL{
			Path: "/api/web/v1/calendar",
		},
		Header: http.Header{"User-Agent": []string{"Googlebot"}, "Host": []string{"testing.com"}},
	}
	location := server.getLocation(req)
	assert.Equal(t, api, location)
}

func TestGetLocationNil(t *testing.T) {
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	locations := make([]*Location, 0)
	locations = append(locations, api)
	server := NewServer("127.0.0.1:8080", &Host{"testing.com", locations})
	req := &http.Request{
		Host: "testing.com",
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

	locations := make([]*Location, 0)
	locations = append(locations, web)
	locations = append(locations, api)
	locations = append(locations, apiWeb)
	locations = append(locations, apiWebHeader)

	server := NewServer("127.0.0.1:8080", &Host{"testing.com", locations})
	server.sortLocations()
	sortedWebHeader := server.hosts[0].locations[0]
	sortedApiWeb := server.hosts[0].locations[1]
	sortedApi := server.hosts[0].locations[2]
	sortedWeb := server.hosts[0].locations[3]
	assert.Equal(t, apiWebHeader, sortedWebHeader)
	assert.Equal(t, apiWeb, sortedApiWeb)
	assert.Equal(t, web, sortedWeb)
	assert.Equal(t, api, sortedApi)
}

func TestServerAddr(t *testing.T) {
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	locations := make([]*Location, 0)
	locations = append(locations, api)
	server := NewServer("127.0.0.1:8080", &Host{"testing.com", locations})
	assert.Equal(t, "127.0.0.1:8080", server.addr)
	server.SetAddr("127.0.0.1:9090")
	assert.Equal(t, "127.0.0.1:9090", server.addr)
}
