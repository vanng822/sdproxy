package sdproxy

import "net/http"

type MatchHeader struct {
	Name    string
	Pattern string
}

type Location struct {
	path     string
	upstream *Upstream
	matches  []*MatchHeader
}

func (loc *Location) Serve(rw http.ResponseWriter, req *http.Request) {
	loc.upstream.Serve(rw, req)
}

func NewLocation(path string, upstream *Upstream, matches ...*MatchHeader) *Location {
	location := &Location{
		path:     path,
		upstream: upstream,
		matches:  matches,
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
	if lbp[i].path < lbp[j].path {
		return true
	}
	if lbp[i].path > lbp[j].path {
		return false
	}
	// JUST simplicity here, assume more match fields should match first
	// when we need sorted for comparing then implement it
	return len(lbp[i].matches) < len(lbp[j].matches)
}
