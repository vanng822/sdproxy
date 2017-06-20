# sdproxy

Simple dev proxy - mimic nginx in dev environment

# Usage

```go
package main

import (
	"log"

	"github.com/vanng822/sdproxy"
)

func main() {
	web := sdproxy.NewLocation("/", sdproxy.NewUpstream(
		sdproxy.NewReverseProxy("127.0.0.1:8090"),
		sdproxy.NewReverseProxy("127.0.0.1:8091")))
	api := sdproxy.NewLocation("/api", sdproxy.NewUpstream(
		sdproxy.NewReverseProxy("127.0.0.1:8092"),
		sdproxy.NewReverseProxy("127.0.0.1:8093")))

	server := sdproxy.NewServer(web, api)
	log.Fatal(server.ListenAndServe("127.0.0.1:9090"))
}
```
