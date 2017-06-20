package sdproxy

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextEndPoint(t *testing.T) {
	r1 := &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8091"
	}}
	r2 := &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8092"
	}}

	up := NewUpstream(r1, r2)
	res := up.nextServer()
	assert.Equal(t, r1, res)
	res = up.nextServer()
	assert.Equal(t, r2, res)
	res = up.nextServer()
	assert.Equal(t, r1, res)
}
