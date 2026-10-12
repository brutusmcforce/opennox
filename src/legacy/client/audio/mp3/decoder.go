// Frame decoding translated from scalar minimp3.
// Upstream authors dedicated this software to the public domain, without warranty.
// CC0 1.0 Universal: https://creativecommons.org/publicdomain/zero/1.0/
package mp3

// MaxSamplesPerFrame is the maximum number of interleaved int16 output values.
const MaxSamplesPerFrame = 1152 * 2

// FrameInfo describes a decoded frame. A scan without a complete frame updates
// only FrameBytes, matching minimp3; callers may clear the other fields first.
type FrameInfo struct {
	FrameBytes  int
	Channels    int
	Hz          int
	Layer       int
	BitrateKbps int
}

// Decoder holds persistent Layer III reservoir and synthesis state. Its zero
// value is ready for use. A Decoder must not be used concurrently.
type Decoder struct {
	state   decoderState
	scratch decodeScratch
}

// Init marks the next frame for resynchronization. Other state is retained until
// that next scan, matching the original reset/seek behavior.
func (d *Decoder) Init() { decoderInit(&d.state) }

// DecodeFrame returns samples per channel and records consumed bytes in info.
// A nil pcm requests metadata only. Otherwise pcm must have room for the frame's
// interleaved samples (MaxSamplesPerFrame suffices for every supported frame).
// This decoder supports Layer III; Layer I/II headers still report metadata.
func (d *Decoder) DecodeFrame(packet []byte, pcm []int16, info *FrameInfo) int {
	if len(packet) > (1<<31-1)-2*maxFreeFormatFrameSize-headerSize {
		panic("mp3: frame input too large")
	}
	dec := &d.state
	i, frameSize := int32(0), int32(0)
	if len(packet) > 4 && dec.header[0] == 0xff && headerCompare(dec.header, headerAt(packet, 0)) {
		frameSize = headerFrameBytes(headerAt(packet, 0), dec.freeFormatBytes) + headerPadding(headerAt(packet, 0))
		if int(frameSize) != len(packet) && (int(frameSize)+headerSize > len(packet) || !headerCompare(headerAt(packet, 0), headerAt(packet, frameSize))) {
			frameSize = 0
		}
	}
	if frameSize == 0 {
		*dec = decoderState{}
		i = findFrame(packet, int32(len(packet)), &dec.freeFormatBytes, &frameSize)
		if frameSize == 0 || int(i+frameSize) > len(packet) {
			info.FrameBytes = int(i)
			return 0
		}
	}
	hdr := headerAt(packet, i)
	dec.header = hdr
	info.FrameBytes = int(i + frameSize)
	info.Channels = 2
	if hdr[3]>>6 == 3 {
		info.Channels = 1
	}
	info.Hz = int(headerSampleRateHz(hdr))
	info.Layer = 4 - int((hdr[1]>>1)&3)
	info.BitrateKbps = int(headerBitrateKbps(hdr))
	samples := int(headerFrameSamples(hdr))
	if pcm == nil {
		return samples
	}
	if len(pcm) < samples*info.Channels {
		panic("mp3: PCM output too small")
	}
	var outer bitReader
	bsInit(&outer, packet[int(i)+headerSize:int(i+frameSize)], int(frameSize)-headerSize)
	if hdr[1]&1 == 0 {
		getBits(&outer, 16)
	}
	if info.Layer != 3 {
		return 0
	}
	// Reuse per-decoder storage to avoid a scratch allocation on every frame.
	// Original C scratch is automatic storage. Go initializes it deterministically;
	// valid decoding does not depend on prior stack contents. The unused QMF lanes
	// are characterized separately in the complete-frame fixture.
	scratch := &d.scratch
	*scratch = decodeScratch{}
	begin := readSideInfo(&outer, &scratch.grInfo, hdr)
	if begin < 0 || outer.pos > outer.limit {
		decoderInit(dec)
		return 0
	}
	success := restoreReservoir(dec, &outer, &scratch.scratchState, begin)
	if success {
		granules := 1
		if hdr[1]&8 != 0 {
			granules = 2
		}
		for gr := 0; gr < granules; gr++ {
			clear(scratch.spectral[:1152])
			decodeLayer3(dec, scratch, scratch.grInfo[gr*info.Channels:], info.Channels)
			synthGranule(dec.qmfState[:], scratch.spectral[:1152], 18, info.Channels, pcm[gr*576*info.Channels:], scratch.synthesis[:])
		}
	}
	saveReservoir(dec, &scratch.scratchState)
	if !success {
		return 0
	}
	return samples
}
