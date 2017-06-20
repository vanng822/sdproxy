package sdproxy

import (
	"net/http"
	"sort"
	"strings"
)

type Server struct {
	locations []*Location
}

// AddLocation will add and sort the paths in reverse natural order
func (s *Server) AddLocation(locations ...*Location) {
	if len(locations) == 0 {
		return
	}
	s.locations = append(s.locations, locations...)
	s.sortLocations()
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, location := range s.locations {
		if strings.HasPrefix(req.URL.RequestURI(), location.path) {
			location.Serve(rw, req)
			return
		}
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Server) sortLocations() {
	sort.Sort(sort.Reverse(LocationByPath(s.locations)))
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

func NewServer(locations ...*Location) *Server {
	server := &Server{}
	server.AddLocation(locations...)
	return server
}
