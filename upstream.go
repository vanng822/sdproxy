package sdproxy

import (
	"net/http"
	"net/http/httputil"
	"sync"
)

type Upstream struct {
	servers     []*httputil.ReverseProxy
	current     int
	currentLock *sync.Mutex
}

func (up *Upstream) AddServer(servers ...*httputil.ReverseProxy) {
	if len(servers) == 0 {
		return
	}
	up.currentLock.Lock()
	defer up.currentLock.Unlock()

	up.servers = append(up.servers, servers...)
}

func (up *Upstream) nextServer() *httputil.ReverseProxy {
	if len(up.servers) == 0 {
		return nil
	}
	up.currentLock.Lock()
	defer up.currentLock.Unlock()
	if up.current >= len(up.servers) {
		up.current = 0
	}
	server := up.servers[up.current]
	up.current++
	return server
}

func (up *Upstream) Serve(rw http.ResponseWriter, req *http.Request) {
	server := up.nextServer()
	if server == nil {
		http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	server.ServeHTTP(rw, req)
}

func NewUpstream(servers ...string) *Upstream {
	up := &Upstream{
		currentLock: &sync.Mutex{},
	}
	if len(servers) > 0 {
		for _, server := range servers {
			up.AddServer(NewReverseProxy(server))
		}
	}
	return up
}
