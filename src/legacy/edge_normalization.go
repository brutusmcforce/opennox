package legacy

func normalizeBorderEdge(index, edge int32) int32 {
	w := int32(edgeAssets[index].Field52)
	h := int32(edgeAssets[index].Field53)
	if w == 3 && h == 3 {
		return edge
	}
	if edge == 0 {
		return 0
	}
	// Keep the ordered signed comparisons even for degenerate dimensions.
	if edge <= w-2 {
		return 1
	}
	if edge == w-1 {
		return 2
	}
	bottom := w + 2*h - 4
	if edge < bottom {
		return 3 + int32((uint32(w)^uint32(edge))&1)
	}
	if edge == bottom {
		return 5
	}
	if edge > 2*(w+h)-6 {
		if edge == 2*(w+h)-5 {
			return 7
		}
		return edge + 2*(6-h-w)
	}
	return 6
}
