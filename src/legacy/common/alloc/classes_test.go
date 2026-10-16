package alloc

import "testing"

type classTestItem struct {
	A, B int
}

// TestDynamicClassGrows: a dynamic class allocates more than its preallocated objects (the client drawables, opennox#715),
// a static one runs out.
func TestDynamicClassGrows(t *testing.T) {
	st := NewClassT("staticTest", classTestItem{}, 2)
	defer st.Free()
	if st.NewObject() == nil || st.NewObject() == nil {
		t.Fatal("static: preallocated objects missing")
	}
	if st.NewObject() != nil {
		t.Fatal("static: expected nil when all are used")
	}

	dyn := NewDynamicClassT("dynamicTest", classTestItem{}, 2)
	defer dyn.Free()
	seen := make(map[*classTestItem]bool)
	var items []*classTestItem
	for i := 0; i < 5; i++ {
		it := dyn.NewObject()
		if it == nil {
			t.Fatalf("dynamic: object %d is nil", i)
		}
		if seen[it] {
			t.Fatalf("dynamic: object %d allocated twice", i)
		}
		seen[it] = true
		it.A = i
		items = append(items, it)
	}
	if n := dyn.Active(); n != 5 {
		t.Fatalf("active: %d, expected 5", n)
	}
	if dyn.Cap() != 2 {
		t.Fatalf("cap: %d, expected 2", dyn.Cap())
	}
	n := 0
	dyn.EachActive(func(it *classTestItem) {
		if !seen[it] {
			t.Errorf("unknown active object %p", it)
		}
		n++
	})
	if n != 5 {
		t.Fatalf("EachActive: %d, expected 5", n)
	}
	for _, it := range items {
		dyn.FreeObjectFirst(it)
	}
	if n := dyn.Active(); n != 0 {
		t.Fatalf("active after free: %d", n)
	}
}
