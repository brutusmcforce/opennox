// Huffman decoding translated from the scalar minimp3 implementation.
// Upstream authors dedicated this software to the public domain, without warranty.
// CC0 1.0 Universal: https://creativecommons.org/publicdomain/zero/1.0/
package mp3

// huffman preserves the original padded scratch-buffer lookahead. bs.limit is
// the logical bit limit; bs.buf capacity describes the actual accessible backing
// storage, including retained bytes beyond that limit. The caller supplies valid
// side-info band tables, up to288 big-value pairs and enough physical lookahead.
// Count1 zero coefficients retain dst values; callers clear their spectral buffer.
func huffman(dst []float32, bs *bitReader, gr *grInfo, scf []float32, granuleLimit int) {
	buf := bs.buf[:cap(bs.buf)]
	next := bs.pos / 8
	cache := (uint32(buf[next])<<24 | uint32(buf[next+1])<<16 | uint32(buf[next+2])<<8 | uint32(buf[next+3])) << uint(bs.pos&7)
	sh := (bs.pos & 7) - 8
	next += 4
	peek := func(n int) uint32 { return cache >> uint(32-n) }
	flush := func(n int) { cache <<= uint(n); sh += n }
	refill := func() {
		for sh >= 0 {
			cache |= uint32(buf[next]) << uint(sh)
			next++
			sh -= 8
		}
	}
	one := float32(0)
	region, big, sfbIndex, scaleIndex, out := 0, int(gr.bigValues), 0, 0, 0
	for big > 0 {
		tab := int(gr.tableSelect[region])
		bands := int(gr.regionCount[region])
		region++
		codebook := huffmanTabs[int(huffmanTabIndex[tab]):]
		linbits := int(huffmanLinbits[tab])
		for {
			np := int(gr.sfbTable[sfbIndex]) / 2
			sfbIndex++
			pairs := big
			if np < pairs {
				pairs = np
			}
			one = scf[scaleIndex]
			scaleIndex++
			for pair := 0; pair < pairs; pair++ {
				w := 5
				leaf := int(codebook[peek(w)])
				for leaf < 0 {
					flush(w)
					w = leaf & 7
					leaf = int(codebook[int(peek(w))-(leaf>>3)])
				}
				flush(leaf >> 8)
				for j := 0; j < 2; j++ {
					lsb := leaf & 15
					if linbits != 0 && lsb == 15 {
						lsb += int(peek(linbits))
						flush(linbits)
						refill()
						v := float32(one * pow43(lsb))
						sign := float32(1)
						if int32(cache) < 0 {
							sign = -1
						}
						dst[out] = float32(v * sign)
					} else {
						dst[out] = float32(pow43Table[16+lsb-16*int(cache>>31)] * one)
					}
					if lsb != 0 {
						flush(1)
					}
					out++
					leaf >>= 4
				}
				refill()
			}
			big -= np
			if big <= 0 {
				break
			}
			bands--
			if bands < 0 {
				break
			}
		}
	}
	np := 1 - big
count1:
	for {
		codebook := huffmanTab32[:]
		if gr.count1Table != 0 {
			codebook = huffmanTab33[:]
		}
		leaf := int(codebook[peek(4)])
		if leaf&8 == 0 {
			leaf = int(codebook[(leaf>>3)+int((cache<<4)>>uint(32-(leaf&3)))])
		}
		flush(leaf & 7)
		if next*8-24+sh > granuleLimit {
			break
		}
		for half := 0; half < 2; half++ {
			np--
			if np == 0 {
				np = int(gr.sfbTable[sfbIndex]) / 2
				sfbIndex++
				if np == 0 {
					break count1
				}
				one = scf[scaleIndex]
				scaleIndex++
			}
			for j := half * 2; j < half*2+2; j++ {
				if leaf&(128>>uint(j)) != 0 {
					v := one
					if int32(cache) < 0 {
						v = -v
					}
					dst[out+j] = v
					flush(1)
				}
			}
		}
		refill()
		out += 4
	}
	bs.pos = granuleLimit
}
