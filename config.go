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
	Addr      string
	Locations []*LocationConfig
}

func NewServerFromConfig(conf *Config) *Server {
	var locations []*Location
	for _, location := range conf.Locations {
		locations = append(locations, NewLocation(location.Path, NewUpstream(location.Servers...), location.Matches...))
	}
	server := NewServer(conf.Addr, locations...)
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
