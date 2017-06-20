package sdproxy

import (
	"net/http"
	"net/http/httputil"
	"sync"
)

type Upstream struct {
	Endpoints   []*httputil.ReverseProxy
	current     int
	currentLock *sync.Mutex
}

func (up *Upstream) AddEndpoint(endpoints ...*httputil.ReverseProxy) {
	if len(endpoints) == 0 {
		return
	}
	up.Endpoints = append(up.Endpoints, endpoints...)
}

func (up *Upstream) nextEndPoint() *httputil.ReverseProxy {
	if len(up.Endpoints) == 0 {
		return nil
	}
	up.currentLock.Lock()
	defer up.currentLock.Unlock()
	if up.current >= len(up.Endpoints) {
		up.current = 0
	}
	endpoint := up.Endpoints[up.current]
	up.current++
	return endpoint
}

func (up *Upstream) Serve(rw http.ResponseWriter, req *http.Request) {
	endpoint := up.nextEndPoint()
	if endpoint == nil {
		http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	endpoint.ServeHTTP(rw, req)
}

func NewUpstream(endpoints ...*httputil.ReverseProxy) *Upstream {
	up := &Upstream{
		currentLock: &sync.Mutex{},
	}
	if len(endpoints) > 0 {
		up.AddEndpoint(endpoints...)
	}
	return up
}
