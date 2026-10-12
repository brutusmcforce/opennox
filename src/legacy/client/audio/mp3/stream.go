// Package mp3 contains Go decoder internals translated from minimp3.
//
// minimp3 is dedicated to the public domain worldwide, to the extent possible
// under law, and is provided without warranty (CC0 1.0 Universal). This
// attribution is retained from minimp3.h.
// Upstream: https://github.com/lieff/minimp3
// License: http://creativecommons.org/publicdomain/zero/1.0/
package mp3

const (
	maxFrameSyncMatches    = 10
	maxFreeFormatFrameSize = 2304
	maxReservoirBytes      = 511
	maxMainDataBytes       = 511 + 2304
	headerSize             = 4
)

// decoderState preserves the C decoder's persistent state and storage extents.
type decoderState struct {
	mdctOverlap     [2][288]float32
	qmfState        [960]float32
	reserv          int32
	freeFormatBytes int32
	header          [4]byte
	reservBuf       [511]byte
}

// scratchState contains the reservoir/reader fields needed by these helpers.
// Other scratch fields belong to later decoder stages.
type scratchState struct {
	bs       bitReader
	mainData [maxMainDataBytes]byte
}

// decoderInit mirrors mp3dec_init: only header[0] is cleared.
func decoderInit(dec *decoderState) {
	dec.header[0] = 0
}

// matchFrame mirrors mp3d_match_frame. packet must contain at least
// mp3Bytes bytes and begin with a structurally valid frame header.
func matchFrame(packet []byte, mp3Bytes, frameBytes int32) bool {
	if mp3Bytes < headerSize || int64(mp3Bytes) > int64(len(packet)) ||
		frameBytes < 0 || frameBytes > maxFreeFormatFrameSize || !headerValid(headerAt(packet, 0)) {
		panic("mp3: unsupported frame-match bounds")
	}
	for i, nmatch := int32(0), int32(0); nmatch < maxFrameSyncMatches; nmatch++ {
		h := headerAt(packet, i)
		i += headerFrameBytes(h, frameBytes) + headerPadding(h)
		if i+headerSize > mp3Bytes {
			return nmatch > 0
		}
		if !headerCompare(headerAt(packet, 0), headerAt(packet, i)) {
			return false
		}
	}
	return true
}

// findFrame mirrors mp3d_find_frame. The scanner checks each read against the
// logical input bounds, so no extra backing padding is required. Output
// pointers are required.
func findFrame(data []byte, mp3Bytes int32, freeFormatBytes, ptrFrameBytes *int32) int32 {
	if freeFormatBytes == nil || ptrFrameBytes == nil || mp3Bytes < 0 || mp3Bytes > (1<<31-1)-2*maxFreeFormatFrameSize-headerSize ||
		int64(mp3Bytes) > int64(len(data)) || *freeFormatBytes < 0 || *freeFormatBytes > maxFreeFormatFrameSize {
		panic("mp3: unsupported frame-finder arguments")
	}
	var i int32
	for i = 0; i < mp3Bytes-headerSize; i++ {
		h := headerAt(data, i)
		if !headerValid(h) {
			continue
		}
		frameBytes := headerFrameBytes(h, *freeFormatBytes)
		frameAndPadding := frameBytes + headerPadding(h)

		for k := int32(headerSize); frameBytes == 0 && k < maxFreeFormatFrameSize && i+2*k < mp3Bytes-headerSize; k++ {
			if !headerCompare(h, headerAt(data, i+k)) {
				continue
			}
			fb := k - headerPadding(h)
			nextfb := fb + headerPadding(headerAt(data, i+k))
			if i+k+nextfb+headerSize > mp3Bytes || !headerCompare(h, headerAt(data, i+k+nextfb)) {
				continue
			}
			frameAndPadding = k
			frameBytes = fb
			*freeFormatBytes = fb
		}
		if (frameBytes != 0 && i+frameAndPadding <= mp3Bytes && matchFrame(data[i:], mp3Bytes-i, frameBytes)) ||
			(i == 0 && frameAndPadding == mp3Bytes) {
			*ptrFrameBytes = frameAndPadding
			return i
		}
		*freeFormatBytes = 0
	}
	*ptrFrameBytes = 0
	return i
}

// saveReservoir mirrors L3_save_reservoir, including unsigned C arithmetic
// before conversion to the signed stored length. Supported positions and
// limits are nonnegative and small enough that the resulting int32 is defined
// on the target C implementation. A nonpositive remainder leaves the buffer
// bytes untouched.
func saveReservoir(dec *decoderState, scratch *scratchState) {
	if scratch.bs.pos < 0 || scratch.bs.limit < 0 || int64(scratch.bs.pos) > (1<<31-1)-7 ||
		scratch.bs.limit > maxMainDataBytes*8 {
		panic("mp3: unsupported reservoir-save reader state")
	}
	pos := (uint32(scratch.bs.pos) + 7) / 8
	remains := int32(uint32(scratch.bs.limit)/8 - pos)
	if remains > maxReservoirBytes {
		pos += uint32(remains - maxReservoirBytes)
		remains = maxReservoirBytes
	}
	if remains > 0 {
		copy(dec.reservBuf[:remains], scratch.mainData[int(pos):int(pos)+int(remains)])
	}
	dec.reserv = remains
}

// restoreReservoir mirrors L3_restore_reservoir. It updates only scratch
// storage and its reader; the caller's reader and decoder state are unchanged.
// Input bounds follow the decoder's frame domain: reservoir and begin are
// 0..511, reader positions are nonnegative, and the current frame contributes
// at most 2304 bytes.
func restoreReservoir(dec *decoderState, bs *bitReader, scratch *scratchState, mainDataBegin int32) bool {
	if dec.reserv < 0 || dec.reserv > maxReservoirBytes || mainDataBegin < 0 || mainDataBegin > maxReservoirBytes ||
		bs.pos < 0 || bs.limit < bs.pos || bs.limit > len(bs.buf)*8 {
		panic("mp3: unsupported reservoir-restore state")
	}
	frameBytes := (bs.limit - bs.pos) / 8
	if frameBytes > 2304 {
		panic("mp3: frame exceeds reservoir scratch capacity")
	}
	bytesHave := int(dec.reserv)
	if int(mainDataBegin) < bytesHave {
		bytesHave = int(mainDataBegin)
	}
	reservoirOffset := int(dec.reserv - mainDataBegin)
	if reservoirOffset < 0 {
		reservoirOffset = 0
	}
	copy(scratch.mainData[:bytesHave], dec.reservBuf[reservoirOffset:reservoirOffset+bytesHave])
	frameStart := bs.pos / 8
	copy(scratch.mainData[bytesHave:bytesHave+frameBytes], bs.buf[frameStart:frameStart+frameBytes])
	bsInit(&scratch.bs, scratch.mainData[:bytesHave+frameBytes], bytesHave+frameBytes)
	return dec.reserv >= mainDataBegin
}

func headerAt(data []byte, off int32) [4]byte {
	i := int(off)
	return [4]byte{data[i], data[i+1], data[i+2], data[i+3]}
}
