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
	var hosts []*Host
    web := sdproxy.NewLocation("/", sdproxy.NewUpstream("127.0.0.1:8090", "127.0.0.1:8091"))
    hosts = append(hosts, &Host{"", web})
	api := sdproxy.NewLocation("/api", sdproxy.NewUpstream("127.0.0.1:8092", "127.0.0.1:8093"))
    hosts = append(hosts, &Host{"", api})

	server := sdproxy.NewServer("127.0.0.1:8181", hosts...)
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
    "hosts": [{
        "hostname": "",
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
    }]
    
}
```
