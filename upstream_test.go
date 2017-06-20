package sdproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextEndPoint(t *testing.T) {
	up := NewUpstream("127.0.0.1:8091", "127.0.0.1:8092")
	r1 := up.servers[0]
	r2 := up.servers[1]
	res := up.nextServer()
	assert.Equal(t, r1, res)
	res = up.nextServer()
	assert.Equal(t, r2, res)
	res = up.nextServer()
	assert.Equal(t, r1, res)
}
