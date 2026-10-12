package protection

import "testing"

func TestInitialize(t *testing.T) {
	sum := uint32(0xaabbccdd)
	if Initialize(nil, 1, 2, 3, &sum) || sum != 0xaabbccdd {
		t.Fatal("allocation failure changed state")
	}
	r := Record{ID: 0xffffffff, Value: 0xffffffff, Next: &Record{}, Prev: &Record{}}
	if !Initialize(&r, 0x80000000, 0x7f800001, 0x12345678, &sum) {
		t.Fatal("initialization failed")
	}
	if r.ID != 0x92345678 || r.Value != 0x6db45679 || r.Next != nil || r.Prev != nil || sum != 0x553bccdc {
		t.Fatalf("record=%+v sum=%08x", r, sum)
	}
}
