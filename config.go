package sdproxy

import (
	"encoding/json"
	"os"
)

type LocationConfig struct {
	Path    string
	Matches []*MatchHeader
	Servers []string
}

type Config struct {
	Addr  string
	Hosts []*HostConfig
}

type HostConfig struct {
	Hostname  string
	Locations []*LocationConfig
}

func NewServerFromConfig(conf *Config) *Server {
	var hosts []*Host
	for _, host := range conf.Hosts {
		var locations []*Location
		for _, location := range host.Locations {
			locations = append(locations, NewLocation(location.Path, NewUpstream(location.Servers...), location.Matches...))
		}
		hosts = append(hosts, &Host{host.Hostname, locations})
	}
	server := NewServer(conf.Addr, hosts...)
	return server
}

func ParseConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
