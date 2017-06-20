package sdproxy

import "net/http"

type Location struct {
	path     string
	upstream *Upstream
}

func (loc *Location) Serve(rw http.ResponseWriter, req *http.Request) {
	loc.upstream.Serve(rw, req)
}

func NewLocation(path string, upstream *Upstream) *Location {
	location := &Location{
		path:     path,
		upstream: upstream,
	}
	return location
}
