package server

import (
	"image"
	"testing"
)

func TestWallDeletePreservesLists(t *testing.T) {
	layouts := []struct {
		name      string
		positions []image.Point
	}{
		{"different_rows", []image.Point{image.Pt(10, 10), image.Pt(12, 12), image.Pt(14, 14), image.Pt(16, 16)}},
		{"same_row", []image.Point{image.Pt(16, 10), image.Pt(10, 10), image.Pt(14, 10), image.Pt(12, 10)}},
		{"same_bucket", []image.Point{image.Pt(10, 10), image.Pt(42, 10), image.Pt(74, 10), image.Pt(106, 10)}},
	}
	for _, layout := range layouts {
		t.Run(layout.name, func(t *testing.T) {
			for _, removed := range []int{0, 1, 2, 3} {
				t.Run([]string{"oldest", "second", "third", "newest"}[removed], func(t *testing.T) {
					s := serverWalls{byPos: make([]*Wall, wallsPerBucket*WallGridSize), indexY: make([]*Wall, WallGridSize)}
					pool := make([]Wall, 4)
					for i := range pool {
						pool[i].Next20 = s.freeList
						s.freeList = &pool[i]
					}
					for _, pos := range layout.positions {
						if s.CreateAtGrid(pos) == nil {
							t.Fatal("wall allocation")
						}
					}
					check := func(want map[image.Point]bool) {
						t.Helper()
						seen := make(map[*Wall]bool)
						for w := s.head; w != nil; w = w.Next20 {
							if seen[w] {
								t.Fatal("cycle in global wall list")
							}
							seen[w] = true
							if !want[w.GridPos()] {
								t.Fatalf("unexpected global wall %v", w.GridPos())
							}
						}
						if len(seen) != len(want) {
							t.Fatalf("global list has %d walls, want %d", len(seen), len(want))
						}
						indexed := make(map[*Wall]bool)
						for i, head := range s.byPos {
							for w := head; w != nil; w = w.NextByPos16 {
								if indexed[w] || !seen[w] {
									t.Fatal("position index cycles, duplicates or includes a non-live wall")
								}
								indexed[w] = true
								at, ok := wallArrayInd(w.GridPos())
								if !ok || int(at) != i {
									t.Fatal("wall in wrong position bucket")
								}
							}
						}
						if len(indexed) != len(seen) {
							t.Fatal("position index lost a live wall")
						}
						byRow := make(map[*Wall]bool)
						for y, head := range s.indexY {
							previous := -1
							for w := head; w != nil; w = w.NextByY24 {
								if byRow[w] || !seen[w] {
									t.Fatal("row index cycles, duplicates or includes a non-live wall")
								}
								byRow[w] = true
								if int(w.Y6) != y || int(w.X5) <= previous {
									t.Fatal("row index has wrong coordinates or order")
								}
								previous = int(w.X5)
							}
						}
						if len(byRow) != len(seen) {
							t.Fatal("row index lost a live wall")
						}
						for pos := range want {
							if w := s.GetWallAtGrid(pos); w == nil || !seen[w] {
								t.Fatalf("lookup lost %v", pos)
							}
						}
						free := make(map[*Wall]bool)
						for w := s.freeList; w != nil; w = w.Next20 {
							if free[w] || seen[w] {
								t.Fatal("free list cycles or overlaps live walls")
							}
							free[w] = true
						}
						if len(free)+len(seen) != len(pool) {
							t.Fatal("wall disappeared from live and free lists")
						}
					}
					want := make(map[image.Point]bool)
					for i, pos := range layout.positions {
						if i != removed {
							want[pos] = true
						}
					}
					s.DeleteAtGrid(layout.positions[removed])
					check(want)
					s.DeleteAtGrid(layout.positions[removed])
					check(want)
					replacement := image.Pt(18, 18)
					if s.CreateAtGrid(replacement) == nil {
						t.Fatal("deleted wall was not reusable")
					}
					want[replacement] = true
					check(want)
				})
			}
		})
	}
}
