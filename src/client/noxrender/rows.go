package noxrender

import "image"

// MergeRows merges spans (one row high) of the same columns in adjacent rows into rectangles,
// so fewer and larger quads draw the same pixels (each as many times as before).
// Spans are given row by row, as the game makes them.
func MergeRows(out []image.Rectangle, spans []image.Rectangle) []image.Rectangle {
	open := 0 // out[open:] may still grow down
	for _, s := range spans {
		merged := false
		for i := open; i < len(out); i++ {
			o := &out[i]
			if o.Max.Y == s.Min.Y && o.Min.X == s.Min.X && o.Max.X == s.Max.X && s.Dy() == 1 {
				o.Max.Y = s.Max.Y
				merged = true
				break
			}
		}
		if !merged {
			// rectangles that end above this row can't grow anymore: move them before the open ones
			for i := open; i < len(out); i++ {
				if out[i].Max.Y < s.Min.Y {
					out[open], out[i] = out[i], out[open]
					open++
				}
			}
			out = append(out, s)
		}
	}
	return out
}
