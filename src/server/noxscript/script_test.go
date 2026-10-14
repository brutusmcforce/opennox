package noxscript

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/brutusmcforce/opennox/noxscript/ns/v4"
)

// stackVM is a VM with a stack of integers and a NoxScript implementation of Random only.
type stackVM struct {
	VM
	stack []uint32
	ns    *randomNS
}

func (s *stackVM) NoxScript() ns.Implementation { return s.ns }
func (s *stackVM) PushU32(v uint32)             { s.stack = append(s.stack, v) }
func (s *stackVM) PushI32(v int32)              { s.PushU32(uint32(v)) }
func (s *stackVM) PopU32() uint32 {
	v := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return v
}
func (s *stackVM) PopI32() int32 { return int32(s.PopU32()) }

type randomNS struct {
	ns.Implementation
	min, max int
}

func (r *randomNS) Random(min, max int) int {
	r.min, r.max = min, max
	return min
}

// Script integers are signed: on 64-bit, Random(-5, 5) once got 4294967291 as its minimum.
func TestRandomIntSigned(t *testing.T) {
	vm := &stackVM{ns: &randomNS{}}
	vm.PushI32(-5)
	vm.PushI32(5)
	nsRandomInt(vm)
	require.Equal(t, -5, vm.ns.min)
	require.Equal(t, 5, vm.ns.max)
	require.Equal(t, int32(-5), vm.PopI32())
}
