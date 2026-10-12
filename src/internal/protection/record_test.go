package protection

import "testing"

func TestRecords(t *testing.T) {
	c := Record{ID: 5, Value: 6}
	b := Record{ID: 3, Value: 4, Next: &c}
	a := Record{ID: 1, Value: 2, Next: &b}
	b.Prev, c.Prev = &a, &b
	if Find(&a, 0x80000000, 0x80000003) != &b || Find(&a, 0, 99) != nil || Find(nil, 0, 1) != nil {
		t.Fatal("ID lookup")
	}
	for _, tc := range []struct {
		index int32
		want  *Record
	}{{-1, nil}, {0, &a}, {1, &b}, {2, &c}, {3, nil}, {2147483647, nil}} {
		if At(&a, tc.index) != tc.want {
			t.Fatalf("index %d", tc.index)
		}
	}
	if Swap(nil, &a) || Swap(&a, nil) || a.ID != 1 || a.Value != 2 {
		t.Fatal("null swap")
	}
	if !Swap(&a, &c) || a.ID != 5 || a.Value != 6 || c.ID != 1 || c.Value != 2 {
		t.Fatal("payload swap")
	}
	if a.Next != &b || a.Prev != nil || b.Next != &c || b.Prev != &a || c.Next != nil || c.Prev != &b {
		t.Fatal("links changed")
	}
	if !Swap(&a, &a) || a.ID != 5 || a.Value != 6 {
		t.Fatal("self swap")
	}
}
