package sdproxy

import (
	"net/http"
	"sort"
	"strings"
)

type Server struct {
	addr      string
	locations []*Location
}

func (s *Server) SetAddr(addr string) {
	s.addr = addr
}

// AddLocation will add and sort the paths in reverse natural order
func (s *Server) AddLocation(locations ...*Location) {
	if len(locations) == 0 {
		return
	}
	s.locations = append(s.locations, locations...)
	s.sortLocations()
}

func (s *Server) getLocation(requestURI string) *Location {
	for _, location := range s.locations {
		if strings.HasPrefix(requestURI, location.path) {
			return location
		}
	}
	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if location := s.getLocation(req.URL.RequestURI()); location != nil {
		location.Serve(rw, req)
		return
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Server) sortLocations() {
	sort.Sort(sort.Reverse(LocationByPath(s.locations)))
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s)
}

func NewServer(addr string, locations ...*Location) *Server {
	server := &Server{
		addr: addr,
	}
	server.AddLocation(locations...)
	return server
}
