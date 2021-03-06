# sdproxy

Simple dev reverse proxy - mimic nginx in dev environment

# Usage

```go
package main

import (
	"log"

	"github.com/vanng822/sdproxy"
)

func main() {
	web := sdproxy.NewLocation("/", sdproxy.NewUpstream("127.0.0.1:8090", "127.0.0.1:8091"))
	api := sdproxy.NewLocation("/api", sdproxy.NewUpstream("127.0.0.1:8092", "127.0.0.1:8093"))

	server := sdproxy.NewServer("127.0.0.1:8181", api, web)
	log.Fatal(server.ListenAndServe())
}
```

OR

```bash
> go install github.com/vanng822/sdproxy/cmd/sdproxy
> sdproxy -c path_to_config.json
```

Configuration example

```json
{
    "addr": "127.0.0.1:8080",
    "locations": [{
        "path": "/",
        "servers": [
            "127.0.0.1:8090",
            "127.0.0.1:8094"
        ]
    }, {
        "path": "/api",
        "servers": [
            "127.0.0.1:8091",
            "127.0.0.1:8092"
        ]
    }]
}
```
