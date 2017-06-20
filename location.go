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

type LocationByPath []*Location

func (lbp LocationByPath) Len() int {
	return len(lbp)
}
func (lbp LocationByPath) Swap(i, j int) {
	lbp[i], lbp[j] = lbp[j], lbp[i]
}
func (lbp LocationByPath) Less(i, j int) bool {
	return lbp[i].path < lbp[j].path
}
