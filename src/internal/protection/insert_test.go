package protection

import "testing"

func TestInsertBefore(t *testing.T) {
	for pos := 0; pos < 3; pos++ {
		old := [3]Record{{ID: 1, Value: 4}, {ID: 2, Value: 5}, {ID: 3, Value: 6}}
		old[0].Next = &old[1]
		old[1].Prev = &old[0]
		old[1].Next = &old[2]
		old[2].Prev = &old[1]
		r := Record{ID: 9, Value: 10}
		head := InsertBefore(&old[0], &old[pos], &r)
		want := []*Record{&old[0], &old[1], &old[2]}
		want = append(want, nil)
		copy(want[pos+1:], want[pos:])
		want[pos] = &r
		p := head
		for i, v := range want {
			if p != v {
				t.Fatalf("position %d node %d", pos, i)
			}
			var prev *Record
			if i > 0 {
				prev = want[i-1]
			}
			if p.Prev != prev {
				t.Fatal("back link")
			}
			p = p.Next
		}
		if p != nil {
			t.Fatal("tail link")
		}
		for i, v := range old {
			if v.ID != uint32(i+1) || v.Value != uint32(i+4) {
				t.Fatal("old payload changed")
			}
		}
		if r.ID != 9 || r.Value != 10 {
			t.Fatal("new payload changed")
		}
	}
}
