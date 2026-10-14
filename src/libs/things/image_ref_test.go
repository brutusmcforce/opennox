package things

import (
	"encoding/json"
	"testing"
)

// ImageRef marshals to the forms UnmarshalJSON reads (it recursed forever for named references).
func TestImageRefJSON(t *testing.T) {
	for _, ref := range []ImageRef{{Ind: 12}, {Ind2: 3, Name: "Foo"}, {Name: "Bar"}, {}} {
		data, err := json.Marshal(ref)
		if err != nil {
			t.Fatal(err)
		}
		var got ImageRef
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("%s: %v", data, err)
		}
		if got != ref {
			t.Fatalf("%s: got %+v, want %+v", data, got, ref)
		}
	}
}
