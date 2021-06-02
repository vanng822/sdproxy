package sdproxy

import (
	"net/http"
	"sort"
	"strings"
)

type Server struct {
	addr  string
	hosts map[string][]*Location
}

func (s *Server) SetAddr(addr string) {
	s.addr = addr
}

// AddLocation will add and sort the paths in reverse natural order
func (s *Server) AddLocation(hosts ...*Host) {
	for _, host := range hosts {
		if len(host.locations) == 0 {
			continue
		}
		if _, exist := s.hosts[host.hostname]; !exist {
			s.hosts[host.hostname] = make([]*Location, 0)
		}
		s.hosts[host.hostname] = append(s.hosts[host.hostname], host.locations...)
	}
	s.sortLocations()
}

func (s *Server) matchHeader(req *http.Request, location *Location) bool {
	for _, header := range location.matches {
		headerValue := req.Header.Get(header.Name)
		if headerValue == "" {
			return false
		}
		if strings.Contains(headerValue, header.Pattern) {
			return true
		}
	}
	return false
}

func (s *Server) getLocation(req *http.Request) *Location {
	for host, locations := range s.hosts {
		if host != req.Host {
			continue
		}
		for _, location := range locations {
			if strings.HasPrefix(req.URL.RequestURI(), location.path) {
				if location.matches == nil {
					return location
				}
				if s.matchHeader(req, location) {
					return location
				}
			}
		}
	}

	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if location := s.getLocation(req); location != nil {
		location.Serve(rw, req)
		return
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Server) sortLocations() {
	for _, locations := range s.hosts {
		sort.Sort(sort.Reverse(LocationByPath(locations)))
	}
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s)
}

func NewServer(addr string, hosts ...*Host) *Server {
	server := &Server{
		addr:  addr,
		hosts: make(map[string][]*Location),
	}
	server.AddLocation(hosts...)
	return server
}
