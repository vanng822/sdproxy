package sdproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	web := NewLocation("/", NewUpstream("127.0.0.1:8090", "127.0.0.1:8094"))
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	apiWeb := NewLocation("/api/web", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))

	server := NewServer("127.0.0.1:8080", web, apiWeb, api)

	location := server.getLocation("/api/web/v1/calendar")
	assert.Equal(t, apiWeb, location)

	location = server.getLocation("/api/v1/calendar")
	assert.Equal(t, api, location)

	location = server.getLocation("/anything/else")
	assert.Equal(t, web, location)
}

func TestGetLocationNil(t *testing.T) {
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	server := NewServer("127.0.0.1:8080", api)
	assert.Nil(t, server.getLocation("/anything/else"))
}

func TestLocationByPath(t *testing.T) {
	web := NewLocation("/", NewUpstream("127.0.0.1:8090", "127.0.0.1:8094"))
	api := NewLocation("/api", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))
	apiWeb := NewLocation("/api/web", NewUpstream("127.0.0.1:8091", "127.0.0.1:8092"))

	server := NewServer("127.0.0.1:8080", web, apiWeb, api)
	server.sortLocations()
	sortedApiWeb := server.locations[0]
	sortedApi := server.locations[1]
	sortedWeb := server.locations[2]

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
