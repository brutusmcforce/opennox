package protection

import (
	"math/rand"
	"testing"
)

func TestRekey(t *testing.T) {
	gen := rand.New(rand.NewSource(0x56f5c0))
	for trial := 0; trial < 1000; trial++ {
		oldKey, newKey := gen.Uint32(), gen.Uint32()
		switch trial % 5 {
		case 0:
			oldKey = 0
		case 1:
			newKey = 0
		case 2:
			newKey = oldKey
		case 3:
			oldKey, newKey = 0xffffffff, 0x80000000
		}
		nodes := make([]Record, trial%51)
		plain := make([][2]uint32, len(nodes))
		var head *Record
		wantSum := ^newKey
		for i := range nodes {
			plain[i] = [2]uint32{gen.Uint32(), gen.Uint32()}
			nodes[i].ID, nodes[i].Value = plain[i][0]^oldKey, plain[i][1]^oldKey
			wantSum ^= plain[i][0] ^ plain[i][1]
			if i == 0 {
				head = &nodes[i]
			} else {
				nodes[i-1].Next, nodes[i].Prev = &nodes[i], &nodes[i-1]
			}
		}
		if got := Rekey(head, oldKey, newKey); got != wantSum {
			t.Fatalf("trial=%d sum=%08x want=%08x", trial, got, wantSum)
		}
		for i := range nodes {
			r := &nodes[i]
			var prev, next *Record
			if i > 0 {
				prev = &nodes[i-1]
			}
			if i+1 < len(nodes) {
				next = &nodes[i+1]
			}
			if r.ID != plain[i][0]^newKey || r.Value != plain[i][1]^newKey || r.Prev != prev || r.Next != next {
				t.Fatalf("trial=%d node=%d: payload or links changed", trial, i)
			}
		}
	}
}
