# sdproxy

Simple dev proxy - mimic nginx in dev environment

# Usage

```go
  api := sdproxy.NewLocation("/api", sdproxy.NewUpstream(&httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8091"
	}}, &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8092"
	}}))

	server := sdproxy.NewServer(api, web)
	log.Fatal(server.ListenAndServe("127.0.0.1:8181"))
  ```
