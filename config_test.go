package sdproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerFromConfig(t *testing.T) {
	conf := Config{
		Addr: "127.0.0.1:8080",
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
	}

	server := NewServerFromConfig(&conf)

	assert.Equal(t, "/api", server.locations[0].path)
	assert.Equal(t, "/", server.locations[1].path)
	assert.Equal(t, "127.0.0.1:8080", server.addr)
}

func TestParseConfig(t *testing.T) {
	conf := ParseConfig("config.json")
	assert.Equal(t, 1, len(conf.Locations))
	assert.Equal(t, "127.0.0.1:8080", conf.Addr)
	location := conf.Locations[0]
	assert.Equal(t, "/api", location.Path)
	assert.Equal(t, []string{"127.0.0.1:8090", "127.0.0.1:8091"}, location.Servers)
}
