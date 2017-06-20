package sdproxy

import (
	"log"
	"net/http"
	"strings"
)

type Server struct {
	locations []*Location
}

// FIFO you are in charge of how it will match
func (s *Server) AddLocation(locations ...*Location) {
	if len(locations) == 0 {
		return
	}
	s.locations = append(s.locations, locations...)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, location := range s.locations {
		log.Println(location.path, req.URL.RequestURI())
		if strings.HasPrefix(req.URL.RequestURI(), location.path) {
			location.Serve(rw, req)
			return
		}
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

func NewServer(locations ...*Location) *Server {
	server := &Server{}
	server.AddLocation(locations...)
	return server
}
