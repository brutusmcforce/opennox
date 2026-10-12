package protection

import "testing"

func TestUnlink(t *testing.T) {
	for remove := 0; remove < 3; remove++ {
		r := [3]Record{{ID: 1, Value: 4}, {ID: 2, Value: 5}, {ID: 3, Value: 6}}
		r[0].Next = &r[1]
		r[1].Prev = &r[0]
		r[1].Next = &r[2]
		r[2].Prev = &r[1]
		head, tail := Unlink(&r[0], &r[2], &r[remove])
		var want []*Record
		for i := range r {
			if i != remove {
				want = append(want, &r[i])
			}
			if r[i].ID != uint32(i+1) || r[i].Value != uint32(i+4) {
				t.Fatal("payload changed")
			}
		}
		if head != want[0] || tail != want[1] || head.Prev != nil || head.Next != tail || tail.Prev != head || tail.Next != nil {
			t.Fatalf("remove %d: broken links", remove)
		}
	}
	r := Record{ID: 1, Value: 2}
	if h, tail := Unlink(&r, &r, &r); h != nil || tail != nil || r.ID != 1 || r.Value != 2 {
		t.Fatal("singleton")
	}
}
