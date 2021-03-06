package sdproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerFromConfig(t *testing.T) {
	conf := Config{
		Addr: "127.0.0.1:8080",
		Hosts: []*HostConfig{
			&HostConfig{
				Hostname: "testing.com",
				Locations: []*LocationConfig{
					&LocationConfig{
						Path: "/",
						Servers: []string{
							"127.0.0.1:8081",
							"127.0.0.1:8082"},
					},
					&LocationConfig{
						Path: "/api",
						Servers: []string{
							"127.0.0.1:8083",
							"127.0.0.1:8084"},
					},
				},
			},
		},
	}

	server := NewServerFromConfig(&conf)

	assert.Equal(t, "/api", server.hosts[0].locations[0].path)
	assert.Equal(t, "/", server.hosts[0].locations[1].path)
	assert.Equal(t, "127.0.0.1:8080", server.addr)
}

func TestParseConfig(t *testing.T) {
	conf := ParseConfig("config.json")
	assert.Equal(t, 1, len(conf.Hosts))
	assert.Equal(t, "127.0.0.1:8080", conf.Addr)
	assert.Equal(t, 2, len(conf.Hosts[0].Locations))
	location := conf.Hosts[0].Locations[0]
	assert.Equal(t, "/api", location.Path)
	assert.Equal(t, []string{"127.0.0.1:8090", "127.0.0.1:8091"}, location.Servers)
	assert.Nil(t, location.Matches)

	location1 := conf.Hosts[0].Locations[1]
	assert.Equal(t, "/", location1.Path)
	assert.Equal(t, []string{"127.0.0.1:8090", "127.0.0.1:8091"}, location1.Servers)
	assert.Equal(t, "User-Agent", location1.Matches[0].Name)
	assert.Equal(t, "Googlebot", location1.Matches[0].Pattern)
}
