//go:build !server

package ail

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/sdl"

	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/log"

	"github.com/brutusmcforce/opennox/legacy/common/alloc/handles"
)

var (
	audioLog     = log.New("audio")
	audioDebug   = os.Getenv("NOX_DEBUG_AUDIO") == "true"
	audioSamples struct {
		sync.RWMutex
		byHandle map[Sample]*audioSample
	}
	audioStreams struct {
		sync.RWMutex
		byHandle map[Stream]*audioStream
	}
	audioDrivers struct {
		sync.RWMutex
		byHandle map[Driver]*audioDriver
	}
	audioTimers struct {
		sync.RWMutex
		byHandle map[Timer]*audioTimer
	}
)

func init() {
	audioSamples.byHandle = make(map[Sample]*audioSample)
	audioStreams.byHandle = make(map[Stream]*audioStream)
	audioDrivers.byHandle = make(map[Driver]*audioDriver)
	audioTimers.byHandle = make(map[Timer]*audioTimer)
}

func (dig Driver) get() *audioDriver {
	if dig == 0 {
		return nil
	}
	handles.AssertValid(uintptr(dig))
	audioDrivers.RLock()
	s := audioDrivers.byHandle[dig]
	audioDrivers.RUnlock()
	return s
}

func (h Stream) get() *audioStream {
	if h == 0 {
		return nil
	}
	handles.AssertValid(uintptr(h))
	audioStreams.RLock()
	s := audioStreams.byHandle[h]
	audioStreams.RUnlock()
	return s
}

func (h Sample) get() *audioSample {
	if h == 0 {
		return nil
	}
	handles.AssertValid(uintptr(h))
	audioSamples.RLock()
	s := audioSamples.byHandle[h]
	audioSamples.RUnlock()
	return s
}

func (h Timer) get() *audioTimer {
	if h == 0 {
		return nil
	}
	handles.AssertValid(uintptr(h))
	audioTimers.RLock()
	s := audioTimers.byHandle[h]
	audioTimers.RUnlock()
	return s
}

// Audio plays through SDL3 audio streams (libs/sdl): each sample and stream has its own, which SDL mixes and converts
// to the device's format. The streams are kept about 100-150 ms ahead by the driver's work, 20 times a second.

type audioDriver struct {
	h          Driver
	mu         sync.Mutex
	sampleHead *audioSample
	streamHead *audioStream
	t          *time.Ticker
}

type audioSampleBuf struct {
	buf []byte
	pos int
}

type audioSample struct {
	h       Sample
	d       *audioDriver
	playing uint32 // atomic

	next      *audioSample
	stream    *sdl.AudioStream
	rate      int32   // sample rate of the queued data (see SetPlaybackRate)
	panL      float32 // gains of the left and right channels for mono data (see SetPan)
	panR      float32
	stereoBuf []int16

	bmu     sync.Mutex
	buffers [2]audioSampleBuf
	current byte
	ready   byte

	stereo bool
	adpcm  bool

	playbackRate uint32
	blockSize    uint32

	eob  func()
	eos  func()
	user any
}

func (s *audioSample) IsPlaying() bool {
	return atomic.LoadUint32(&s.playing) != 0
}

func (s *audioSample) Play() {
	atomic.StoreUint32(&s.playing, 1)
}

func (s *audioSample) Stop() {
	atomic.StoreUint32(&s.playing, 0)
}

func (s *audioSample) EOS() {
	s.Stop()
	if s.eos != nil {
		s.d.mu.Unlock()
		s.eos()
		s.d.mu.Lock()
	}
}

func (s *audioSample) EOB() {
	s.bmu.Lock()
	s.ready++
	if s.current == 0 {
		s.current = 1
	} else {
		s.current = 0
	}
	s.bmu.Unlock()
	if s.eob != nil {
		s.d.mu.Unlock()
		s.eob()
		s.d.mu.Lock()
	}
}

type audioStream struct {
	h Stream
	d *audioDriver

	next   *audioStream
	stream *sdl.AudioStream

	playing bool
	paused  bool

	dec audioDecoder
}

func (s *audioStream) Channels() int {
	return s.dec.Channels()
}

type audioTimer struct {
	h    Timer
	t    *time.Ticker
	dt   time.Duration
	f    func(u uint32)
	user uint32
}

// sampleQueueAhead is how much sound effect data is kept queued (the driver works 20 times a second).
const sampleQueueAhead = 150 * time.Millisecond

// streamQueueAhead is how much music and speech is kept queued.
const streamQueueAhead = 200 * time.Millisecond

func queueBytes(rate, channels int, d time.Duration) int {
	return int(int64(rate) * int64(channels) * 2 * int64(d) / int64(time.Second))
}

func (dig Driver) AllocateSample() Sample {
	if env.IsE2E() {
		h := handles.New()
		return Sample(h)
	}
	if audioDebug {
		audioLog.Println("AIL_allocate_sample_handle")
	}
	d := dig.get()
	if d == nil {
		return 0
	}

	// sound effects are mono or stereo; mono is turned into stereo with the pan applied, see queue
	stream, err := sdl.OpenAudioStream(sdl.AudioSpec{Format: sdl.AUDIO_S16, Channels: 2, Freq: 22050})
	if err != nil {
		audioLog.Println("cannot open audio stream:", err)
		return 0
	}
	s := &audioSample{
		h:      Sample(handles.New()),
		d:      d,
		stream: stream,
		rate:   22050,
	}
	s.setPan(63)

	d.mu.Lock()
	s.next = d.sampleHead
	d.sampleHead = s
	d.mu.Unlock()

	audioSamples.Lock()
	audioSamples.byHandle[s.h] = s
	audioSamples.Unlock()
	return s.h
}

func (s Sample) Release() {
	if audioDebug {
		audioLog.Println("AIL_release_sample_handle")
	}
	handles.AssertValid(uintptr(s))
	audioSamples.Lock()
	v := audioSamples.byHandle[s]
	delete(audioSamples.byHandle, s)
	audioSamples.Unlock()
	if v == nil {
		return
	}
	d := v.d
	d.mu.Lock()
	defer d.mu.Unlock()
	// the driver must not use it anymore; buffers can only be deleted once the source let go of them
	for p := &d.sampleHead; *p != nil; p = &(*p).next {
		if *p == v {
			*p = v.next
			break
		}
	}
	atomic.StoreUint32(&v.playing, 0)
	v.stream.Destroy()
}

func openAudio(path string) (audioDecoder, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	default:
		return nil, fmt.Errorf("unsupported audio file format: %s", path)
	case ".wav":
		r, err := openWav(path)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}

func (dig Driver) OpenStream(name string) Stream {
	if audioDebug {
		audioLog.Println("AIL_open_stream")
	}
	if env.IsE2E() {
		h := handles.New()
		return Stream(h)
	}
	d := dig.get()
	if d == nil {
		return 0
	}

	dec, err := openAudio(name)
	if err != nil {
		audioLog.Println(err)
		return 0
	}

	stream, err := sdl.OpenAudioStream(sdl.AudioSpec{Format: sdl.AUDIO_S16, Channels: int32(dec.Channels()), Freq: int32(dec.SampleRate())})
	if err != nil {
		audioLog.Println("cannot open audio stream:", err)
		dec.Close()
		return 0
	}
	s := &audioStream{
		d:      d,
		h:      Stream(handles.New()),
		dec:    dec,
		stream: stream,
	}

	d.mu.Lock()
	s.next = d.streamHead
	d.streamHead = s
	d.mu.Unlock()

	audioStreams.Lock()
	audioStreams.byHandle[s.h] = s
	audioStreams.Unlock()
	return s.h
}

func (h Stream) Close() error {
	if audioDebug {
		audioLog.Println("AIL_close_stream")
	}
	if env.IsE2E() {
		handles.AssertValid(uintptr(h))
		return nil
	}
	s := h.get()
	if s == nil {
		return nil
	}
	d := s.d
	d.mu.Lock()
	s.playing = false
	s.dec.Close()
	// every stream has its own SDL stream and logical device: free them, or they pile up
	s.stream.Destroy()
	for p := &d.streamHead; *p != nil; p = &(*p).next {
		if *p == s {
			*p = s.next
			break
		}
	}
	d.mu.Unlock()
	audioStreams.Lock()
	delete(audioStreams.byHandle, h)
	audioStreams.Unlock()
	return nil
}

func RegisterTimer(f func(u uint32)) Timer {
	if audioDebug {
		audioLog.Println("AIL_register_timer")
	}
	if env.IsE2E() {
		h := handles.New()
		return Timer(h)
	}
	h := Timer(handles.New())
	t := &audioTimer{h: h, f: f}

	audioTimers.Lock()
	audioTimers.byHandle[h] = t
	audioTimers.Unlock()
	return h
}

func (h Timer) Release() {
	if audioDebug {
		audioLog.Println("AIL_release_timer_handle")
	}
	handles.AssertValid(uintptr(h))
	if t := audioTimers.byHandle[h]; t != nil {
		if t.t != nil {
			t.t.Stop()
		}
		delete(audioTimers.byHandle, h)
	}
}

func (h Sample) End() {
	if audioDebug {
		audioLog.Println("AIL_end_sample")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}

	s.d.mu.Lock()
	defer s.d.mu.Unlock()
	if s.IsPlaying() {
		s.EOB()
		s.EOS()
	}
}

func (h Sample) Init() {
	if audioDebug {
		audioLog.Println("AIL_init_sample")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}

	s.d.mu.Lock()
	defer s.d.mu.Unlock()

	s.bmu.Lock()
	s.current = 0
	s.ready = 2
	s.bmu.Unlock()

	s.stream.SetGain(1)
	s.setPan(63)

	s.Stop()
	// drop what's still queued, like stopping the OpenAL source did
	s.stream.Pause()
	s.stream.Clear()
}

func LastError() string {
	if audioDebug {
		audioLog.Println("AIL_last_error")
	}
	return ""
}

func (h Sample) LoadBuffer(num uint32, buf []byte) {
	if audioDebug {
		audioLog.Printf("AIL_load_sample_buffer: [%d]", len(buf))
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.bmu.Lock()
	defer s.bmu.Unlock()
	sb := &s.buffers[num]
	sb.buf = buf
	sb.pos = 0
	s.Play()
}

func (h Stream) Pause(pause bool) {
	if audioDebug {
		audioLog.Println("AIL_pause_stream")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.d.mu.Lock()
	defer s.d.mu.Unlock()
	if s.paused == pause {
		return
	}
	s.paused = pause
	if pause {
		s.stream.Pause()
	} else if s.playing {
		s.stream.Resume()
	}
}

func (h Sample) RegisterEOBCallback(f func()) {
	if audioDebug {
		audioLog.Println("AIL_register_EOB_callback")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.eob = f
}

func (h Sample) RegisterEOSCallback(f func()) {
	if audioDebug {
		audioLog.Println("AIL_register_EOS_callback")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.eos = f
}

func (h Sample) BufferReady() int {
	if audioDebug {
		audioLog.Println("AIL_sample_buffer_ready")
	}
	if env.IsE2E() {
		return -1
	}
	s := h.get()
	if s == nil {
		return -1
	}
	s.bmu.Lock()
	defer s.bmu.Unlock()
	if s.ready == 0 {
		return -1
	}
	if s.ready == 2 {
		s.ready--
		return 0
	}
	s.ready--
	if s.current == 0 {
		return 1
	}
	return 0
}

func (h Sample) UserData() any {
	if audioDebug {
		audioLog.Println("AIL_sample_user_data")
	}
	if env.IsE2E() {
		return nil
	}
	s := h.get()
	if s == nil {
		return nil
	}
	return s.user
}

func Serve() {
	if audioDebug {
		audioLog.Println("AIL_serve")
	}
	for _, d := range audioDrivers.byHandle {
		select {
		case <-d.t.C:
			d.doWork()
		default:
		}
	}
	for _, t := range audioTimers.byHandle {
		select {
		case <-t.t.C:
			t.f(t.user)
			t.t.Reset(t.dt)
		default:
		}
	}
}

func (h Sample) SetADPCMBlockSize(block uint32) {
	if audioDebug {
		audioLog.Println("AIL_set_sample_adpcm_block_size")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.blockSize = block
}

func (h Sample) SetPan(pan int) {
	if audioDebug {
		audioLog.Println("AIL_set_sample_pan")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.d.mu.Lock()
	defer s.d.mu.Unlock()
	s.setPan(pan)
}

// setPan sets the left and right gains for mono data with a constant power pan law, like OpenAL panned a mono source
// placed on a unit circle in front of the listener (pan 0 is left, 63 the center, 127 right).
// It applies to the data queued from now on.
func (s *audioSample) setPan(pan int) {
	x := min(max(float64(pan-63)/64.0, -1), 1)
	angle := (x + 1) * math.Pi / 4
	s.panL = float32(math.Cos(angle))
	s.panR = float32(math.Sin(angle))
}

func (h Sample) SetPlaybackRate(rate int) {
	if audioDebug {
		audioLog.Println("AIL_set_sample_playback_rate")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.playbackRate = uint32(rate)
}

func (h Sample) SetType(format int32, flags uint32) {
	if audioDebug {
		audioLog.Println("AIL_set_sample_type")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	if flags != 0 {
		panic("abort")
	}
	s.stereo = format&2 != 0
	s.adpcm = format&4 != 0
}

func (h Sample) SetUserData(value any) {
	if audioDebug {
		audioLog.Println("AIL_set_sample_user_data")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.user = value
}

func (h Sample) SetVolume(volume int) {
	if audioDebug {
		audioLog.Println("AIL_set_sample_volume")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.stream.SetGain(float32(volume) / 127.0)
}

func (h Stream) SetPosition(offset int) {
	if audioDebug {
		audioLog.Println("AIL_set_stream_position")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.d.mu.Lock()
	defer s.d.mu.Unlock()
	s.dec.Seek(offset)
}

func (h Stream) SetVolume(volume int) {
	if audioDebug {
		audioLog.Println("AIL_set_stream_volume", int(volume))
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.stream.SetGain(float32(volume) / 127.0)
}

func (h Timer) SetFrequency(hertz uint) {
	if audioDebug {
		audioLog.Println("AIL_set_timer_frequency")
	}
	if env.IsE2E() {
		return
	}
	t := h.get()
	if t == nil {
		return
	}
	t.dt = time.Duration(1000.0/float32(hertz)) * time.Millisecond
}

func (h Stream) Start() {
	if audioDebug {
		audioLog.Println("AIL_start_stream")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.playing = true
}

func (h Timer) Start() {
	if audioDebug {
		audioLog.Println("AIL_start_timer")
	}
	if env.IsE2E() {
		return
	}
	t := h.get()
	if t == nil {
		return
	}
	if t.t != nil {
		t.t.Stop()
	}
	t.t = time.NewTicker(t.dt)
}

func Startup() int {
	if audioDebug {
		audioLog.Println("AIL_startup")
	}
	return -1
}

func Shutdown() {
	if audioDebug {
		audioLog.Println("AIL_shutdown")
	}
	if env.IsE2E() {
		return
	}
	audioTimers.Lock()
	for _, t := range audioTimers.byHandle {
		if t.t != nil {
			t.t.Stop()
		}
	}
	audioTimers.Unlock()
	audioDrivers.Lock()
	for _, d := range audioDrivers.byHandle {
		if d.t != nil {
			d.t.Stop()
		}
	}
	audioDrivers.Unlock()
}

func (h Sample) Stop() {
	if audioDebug {
		audioLog.Println("AIL_stop_sample")
	}
	if env.IsE2E() {
		return
	}
	s := h.get()
	if s == nil {
		return
	}
	s.Stop() // TODO: anything else here?
}

func (h Timer) Stop() {
	if audioDebug {
		audioLog.Println("AIL_stop_timer")
	}
	if env.IsE2E() {
		return
	}
	t := h.get()
	if t == nil {
		return
	}
	if t.t != nil {
		t.t.Stop()
	}
}

func (h Stream) Position() int {
	if audioDebug {
		audioLog.Println("AIL_stream_position")
	}
	if env.IsE2E() {
		return -1
	}
	s := h.get()
	if s == nil {
		return -1
	}
	return s.dec.Position()
}

func (h Stream) Status() int {
	if env.IsE2E() {
		return 2
	}
	s := h.get()
	if s == nil {
		return 2
	}
	s.d.mu.Lock()
	defer s.d.mu.Unlock()
	switch {
	case s.playing && s.paused:
		return 8 // SMP_STOPPED
	case s.playing:
		return 4 // SMP_PLAYING
	}
	return 2 // SMP_DONE
}

func (dig Driver) Close() error {
	if audioDebug {
		audioLog.Println("AIL_waveOutClose")
	}
	if env.IsE2E() {
		return nil
	}
	d := dig.get()
	if d == nil {
		return nil
	}
	if d.t != nil {
		d.t.Stop()
	}
	return nil
}

// queue decodes an ADPCM block and queues it, in stereo with the pan applied to mono data.
func (s *audioSample) queue(data []byte) {
	if !s.adpcm {
		panic("abort")
	}
	var decoded []int16
	if s.stereo {
		decoded = decodeADPCMStereo(decoded, data)
	} else {
		decoded = decodeADPCMMono(decoded, data)
		out := s.stereoBuf[:0]
		for _, v := range decoded {
			out = append(out, int16(float32(v)*s.panL), int16(float32(v)*s.panR))
		}
		s.stereoBuf = out
		decoded = out
	}
	if rate := int32(s.playbackRate); rate > 0 && rate != s.rate {
		if err := s.stream.SetFormat(sdl.AudioSpec{Format: sdl.AUDIO_S16, Channels: 2, Freq: rate}); err != nil {
			audioLog.Println("cannot set the sample rate:", err)
			return
		}
		s.rate = rate
	}
	if err := s.stream.PutInt16(decoded); err != nil {
		audioLog.Println("cannot queue audio:", err)
	}
}

func (s *audioSample) work() {
	if !s.IsPlaying() {
		return
	}

	s.bmu.Lock()
	if s.ready == 2 {
		s.bmu.Unlock()
		s.EOS()
		return
	}

	buf := &s.buffers[s.current]
	if len(buf.buf) == 0 {
		s.bmu.Unlock()
		s.EOB()
		s.EOS()
		return
	}
	defer s.bmu.Unlock()

	ahead := queueBytes(int(s.rate), 2, sampleQueueAhead)
	for s.stream.Queued() < ahead {
		s.queue(buf.buf[buf.pos : buf.pos+int(s.blockSize)])
		buf.pos += int(s.blockSize)

		if s.IsPlaying() && s.stream.Paused() {
			s.stream.Resume()
		}

		if buf.pos >= len(buf.buf) {
			s.bmu.Unlock()
			s.EOB()
			s.bmu.Lock()
			buf = &s.buffers[s.current]
			if len(buf.buf) == 0 {
				break
			}
		}
	}
}

func (s *audioStream) work() {
	if !s.playing || s.paused {
		return
	}
	var buffer [16 * 1024 * 3]int16
	channels := s.Channels()
	sampleRate := s.dec.SampleRate()
	// decode in chunks of 100 ms
	minSamples := sampleRate * channels / 10
	ahead := queueBytes(sampleRate, channels, streamQueueAhead)

	for s.stream.Queued() < ahead {
		offset := 0
		for offset < minSamples {
			samples, ok := s.dec.Decode(buffer[offset : minSamples*2])
			if !ok {
				s.playing = false
			}
			if samples == 0 {
				break
			}
			offset += int(samples)
		}
		if offset == 0 {
			break
		}
		if err := s.stream.PutInt16(buffer[:offset]); err != nil {
			audioLog.Println("cannot queue audio:", err)
			return
		}
		if s.stream.Paused() {
			s.stream.Resume()
		}
	}
}

func (d *audioDriver) doWork() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for s := d.sampleHead; s != nil; s = s.next {
		s.work()
	}
	for s := d.streamHead; s != nil; s = s.next {
		s.work()
	}
}

func WaveOutOpen() Driver {
	if audioDebug {
		audioLog.Println("AIL_waveOutOpen")
	}
	if env.IsE2E() {
		return 0
	}
	h := Driver(handles.New())

	if err := sdl.InitSubSystem(sdl.INIT_AUDIO); err != nil {
		audioLog.Println("cannot open the audio device:", err)
		return 0
	}
	// check that a stream can be opened: without a device, SDL fails here
	test, err := sdl.OpenAudioStream(sdl.AudioSpec{Format: sdl.AUDIO_S16, Channels: 2, Freq: 22050})
	if err != nil {
		audioLog.Println("cannot open the audio device:", err)
		return 0
	}
	test.Destroy()
	audioLog.Println("device ok")

	d := &audioDriver{h: h}
	d.t = time.NewTicker(time.Second / 20)

	audioDrivers.Lock()
	audioDrivers.byHandle[h] = d
	audioDrivers.Unlock()

	return h
}

// PCMStream plays 16-bit PCM data that the caller queues itself (the movie player).
type PCMStream struct {
	s          *sdl.AudioStream
	frameBytes int
	put        atomic.Int64 // sample frames queued so far
}

// OpenPCMStream opens a stream for 16-bit PCM data. It returns nil without audio. Call Start to play.
func (dig Driver) OpenPCMStream(rate, channels int) *PCMStream {
	if env.IsE2E() || dig.get() == nil {
		return nil
	}
	s, err := sdl.OpenAudioStream(sdl.AudioSpec{Format: sdl.AUDIO_S16, Channels: int32(channels), Freq: int32(rate)})
	if err != nil {
		audioLog.Println("cannot open audio stream:", err)
		return nil
	}
	return &PCMStream{s: s, frameBytes: 2 * channels}
}

// QueueStereo16 queues stereo sample frames.
func (p *PCMStream) QueueStereo16(frames [][2]int16) {
	if len(frames) == 0 {
		return
	}
	if err := p.s.PutInt16(unsafe.Slice(&frames[0][0], len(frames)*2)); err != nil {
		audioLog.Println("cannot queue audio:", err)
		return
	}
	p.put.Add(int64(len(frames)))
}

// Played returns the number of sample frames played so far.
func (p *PCMStream) Played() int {
	return max(int(p.put.Load())-p.s.Queued()/p.frameBytes, 0)
}

// Queued returns the number of sample frames not played yet.
func (p *PCMStream) Queued() int {
	return p.s.Queued() / p.frameBytes
}

func (p *PCMStream) SetGain(gain float32) { p.s.SetGain(gain) }
func (p *PCMStream) Start()               { p.s.Resume() }
func (p *PCMStream) Pause()               { p.s.Pause() }
func (p *PCMStream) Close()               { p.s.Destroy() }
