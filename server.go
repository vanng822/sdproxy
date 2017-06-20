package sdproxy

import (
	"log"
	"net/http"
	"strings"
)

type SDProxy struct {
	locations []*Location
}

// FIFO you are in charge of how it will match
func (sdp *SDProxy) AddLocation(locations ...*Location) {
	if len(locations) == 0 {
		return
	}
	sdp.locations = append(sdp.locations, locations...)
}

func (sdp *SDProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, location := range sdp.locations {
		log.Println(location.path, req.URL.RequestURI())
		if strings.HasPrefix(req.URL.RequestURI(), location.path) {
			location.Serve(rw, req)
			return
		}
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (sdp *SDProxy) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, sdp)
}

func NewServer(locations ...*Location) *SDProxy {
	server := &SDProxy{}
	server.AddLocation(locations...)
	return server
}
