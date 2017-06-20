package main

import (
	"flag"
	"log"

	"github.com/vanng822/sdproxy"
)

func main() {
	var (
		config string
	)
	const (
		usage = "sdproxy -h"
	)

	flag.StringVar(&config, "c", "config.json", usage)
	conf := sdproxy.ParseConfig(config)
	server := sdproxy.NewServerFromConfig(conf)
	log.Fatal(server.ListenAndServe(conf.Addr))
}
