// Package openal is a minimal OpenAL binding without cgo: the library is loaded at runtime (see dynlib).
//
// It covers what OpenNox uses, with the same API as the go-openal package it replaces.
// Calls with float arguments use the vector forms (alSourcefv), since purego.SyscallN can't pass floats.
package openal

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"

	"github.com/brutusmcforce/opennox/libs/dynlib"
)

func libNames() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"OpenAL32.dll", "soft_oal.dll"}
	case "darwin":
		return []string{"libopenal.1.dylib", "libopenal.dylib", "/System/Library/Frameworks/OpenAL.framework/OpenAL"}
	}
	return []string{"libopenal.so.1", "libopenal.so"}
}

var fn struct {
	alcOpenDevice          uintptr
	alcCloseDevice         uintptr
	alcCreateContext       uintptr
	alcDestroyContext      uintptr
	alcMakeContextCurrent  uintptr
	alGetError             uintptr
	alGetString            uintptr
	alGenSources           uintptr
	alDeleteSources        uintptr
	alGenBuffers           uintptr
	alDeleteBuffers        uintptr
	alBufferData           uintptr
	alSourcei              uintptr
	alSourcefv             uintptr
	alGetSourcei           uintptr
	alSourcePlay           uintptr
	alSourceStop           uintptr
	alSourcePause          uintptr
	alSourceQueueBuffers   uintptr
	alSourceUnqueueBuffers uintptr
	alListenerfv           uintptr
	alGetBufferi           uintptr
}

var (
	loadOnce sync.Once
	loadErr  error
)

// Load loads the OpenAL library. OpenDevice calls it; it returns nil if the library can't be loaded.
func Load() error {
	loadOnce.Do(func() {
		lib, err := dynlib.Open(libNames()...)
		if err != nil {
			loadErr = err
			return
		}
		loadErr = lib.Syms(map[string]*uintptr{
			"alcOpenDevice":          &fn.alcOpenDevice,
			"alcCloseDevice":         &fn.alcCloseDevice,
			"alcCreateContext":       &fn.alcCreateContext,
			"alcDestroyContext":      &fn.alcDestroyContext,
			"alcMakeContextCurrent":  &fn.alcMakeContextCurrent,
			"alGetError":             &fn.alGetError,
			"alGetString":            &fn.alGetString,
			"alGenSources":           &fn.alGenSources,
			"alDeleteSources":        &fn.alDeleteSources,
			"alGenBuffers":           &fn.alGenBuffers,
			"alDeleteBuffers":        &fn.alDeleteBuffers,
			"alBufferData":           &fn.alBufferData,
			"alSourcei":              &fn.alSourcei,
			"alSourcefv":             &fn.alSourcefv,
			"alGetSourcei":           &fn.alGetSourcei,
			"alSourcePlay":           &fn.alSourcePlay,
			"alSourceStop":           &fn.alSourceStop,
			"alSourcePause":          &fn.alSourcePause,
			"alSourceQueueBuffers":   &fn.alSourceQueueBuffers,
			"alSourceUnqueueBuffers": &fn.alSourceUnqueueBuffers,
			"alListenerfv":           &fn.alListenerfv,
			"alGetBufferi":           &fn.alGetBufferi,
		})
	})
	return loadErr
}

func call(f uintptr, args ...uintptr) uintptr {
	r, _, _ := purego.SyscallN(f, args...)
	return r
}

var (
	ErrInvalidName      = errors.New("openal: invalid name")
	ErrInvalidEnum      = errors.New("openal: invalid enum")
	ErrInvalidValue     = errors.New("openal: invalid value")
	ErrInvalidOperation = errors.New("openal: invalid operation")
)

type ErrorCode uint32

func (e ErrorCode) Error() string {
	return fmt.Sprintf("openal: error code %x", uint32(e))
}

// Err returns the last OpenAL error and clears it. It returns the loading error if the library can't be loaded.
func Err() error {
	if err := Load(); err != nil {
		return err
	}
	switch code := uint32(call(fn.alGetError)); code {
	case 0x0000:
		return nil
	case 0xA001:
		return ErrInvalidName
	case 0xA002:
		return ErrInvalidEnum
	case 0xA003:
		return ErrInvalidValue
	case 0xA004:
		return ErrInvalidOperation
	default:
		return ErrorCode(code)
	}
}

const alVendor = 0xB001

// GetVendor returns the name of the OpenAL implementation.
func GetVendor() string {
	if Load() != nil {
		return ""
	}
	return dynlib.GoString(call(fn.alGetString, alVendor))
}

// Device is an audio device (ALCdevice).
type Device struct {
	handle uintptr
}

// OpenDevice opens an audio device ("" is the default one). It returns nil if there's no device, or no OpenAL.
func OpenDevice(name string) *Device {
	if Load() != nil {
		return nil
	}
	var p uintptr
	if name != "" {
		p = uintptr(unsafe.Pointer(dynlib.CString(name)))
	}
	h := call(fn.alcOpenDevice, p)
	if h == 0 {
		return nil
	}
	return &Device{handle: h}
}

func (d *Device) CloseDevice() bool {
	return call(fn.alcCloseDevice, d.handle)&0xff != 0
}

// CreateContext creates a context on the device, or returns nil.
func (d *Device) CreateContext() *Context {
	h := call(fn.alcCreateContext, d.handle, 0)
	if h == 0 {
		return nil
	}
	return &Context{handle: h}
}

// Context is an OpenAL context (ALCcontext).
type Context struct {
	handle uintptr
}

// Activate makes the context current.
func (c *Context) Activate() bool {
	return call(fn.alcMakeContextCurrent, c.handle)&0xff != 0
}

func (c *Context) Destroy() {
	call(fn.alcDestroyContext, c.handle)
}

// Vector is a 3D vector.
type Vector [3]float32

// Parameters.
const (
	AlPitch            = 0x1003
	AlPosition         = 0x1004
	AlBuffer           = 0x1009
	AlGain             = 0x100A
	AlSourceState      = 0x1010
	AlBuffersQueued    = 0x1015
	AlBuffersProcessed = 0x1016
)

// State is a source state.
type State int32

const (
	Initial State = 0x1011
	Playing State = 0x1012
	Paused  State = 0x1013
	Stopped State = 0x1014
)

func (s State) String() string {
	switch s {
	case Initial:
		return "Initial"
	case Playing:
		return "Playing"
	case Paused:
		return "Paused"
	case Stopped:
		return "Stopped"
	}
	return fmt.Sprintf("State(%#x)", int32(s))
}

// Format is a buffer data format.
type Format int32

const (
	FormatMono8    Format = 0x1100
	FormatMono16   Format = 0x1101
	FormatStereo8  Format = 0x1102
	FormatStereo16 Format = 0x1103
)

// Buffer is an OpenAL buffer name.
type Buffer uint32

// Buffers is a list of buffers.
type Buffers []Buffer

// NewBuffers creates n buffers.
func NewBuffers(n int) Buffers {
	b := make(Buffers, n)
	call(fn.alGenBuffers, uintptr(n), uintptr(unsafe.Pointer(&b[0])))
	return b
}

// NewBuffer creates a buffer.
func NewBuffer() Buffer {
	var b Buffer
	call(fn.alGenBuffers, 1, uintptr(unsafe.Pointer(&b)))
	return b
}

func (b Buffers) Delete() {
	if len(b) == 0 {
		return
	}
	call(fn.alDeleteBuffers, uintptr(len(b)), uintptr(unsafe.Pointer(&b[0])))
}

func (b Buffer) Delete() {
	call(fn.alDeleteBuffers, 1, uintptr(unsafe.Pointer(&b)))
}

func (b Buffer) SetData(format Format, data []byte, frequency int32) {
	call(fn.alBufferData, uintptr(b), uintptr(format), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(frequency))
}

func (b Buffer) SetDataInt16(format Format, data []int16, frequency int32) {
	call(fn.alBufferData, uintptr(b), uintptr(format), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)*2), uintptr(frequency))
}

func (b Buffer) SetDataStereo16(data [][2]int16, frequency int32) {
	call(fn.alBufferData, uintptr(b), uintptr(FormatStereo16), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)*4), uintptr(frequency))
}

// Source is an OpenAL source name.
type Source uint32

// NewSource creates a source.
func NewSource() Source {
	var s Source
	call(fn.alGenSources, 1, uintptr(unsafe.Pointer(&s)))
	return s
}

func (s Source) Delete() {
	call(fn.alDeleteSources, 1, uintptr(unsafe.Pointer(&s)))
}

func (s Source) Seti(param int32, value int32) {
	call(fn.alSourcei, uintptr(s), uintptr(param), uintptr(value))
}

func (s Source) Geti(param int32) int32 {
	var v int32
	call(fn.alGetSourcei, uintptr(s), uintptr(param), uintptr(unsafe.Pointer(&v)))
	return v
}

func (s Source) Setf(param int32, value float32) {
	call(fn.alSourcefv, uintptr(s), uintptr(param), uintptr(unsafe.Pointer(&value)))
}

func (s Source) Setfv(param int32, values []float32) {
	call(fn.alSourcefv, uintptr(s), uintptr(param), uintptr(unsafe.Pointer(&values[0])))
}

func (s Source) SetGain(gain float32) {
	s.Setf(AlGain, gain)
}

func (s Source) SetPosition(v *Vector) {
	call(fn.alSourcefv, uintptr(s), AlPosition, uintptr(unsafe.Pointer(v)))
}

func (s Source) Play()  { call(fn.alSourcePlay, uintptr(s)) }
func (s Source) Stop()  { call(fn.alSourceStop, uintptr(s)) }
func (s Source) Pause() { call(fn.alSourcePause, uintptr(s)) }

func (s Source) State() State {
	return State(s.Geti(AlSourceState))
}

func (s Source) QueueBuffer(b Buffer) {
	call(fn.alSourceQueueBuffers, uintptr(s), 1, uintptr(unsafe.Pointer(&b)))
}

func (s Source) UnqueueBuffers(b Buffers) {
	if len(b) == 0 {
		return
	}
	call(fn.alSourceUnqueueBuffers, uintptr(s), uintptr(len(b)), uintptr(unsafe.Pointer(&b[0])))
}

// Listener is the OpenAL listener.
type Listener struct{}

func (Listener) Setf(param int32, value float32) {
	call(fn.alListenerfv, uintptr(param), uintptr(unsafe.Pointer(&value)))
}

func (Listener) Setfv(param int32, values []float32) {
	call(fn.alListenerfv, uintptr(param), uintptr(unsafe.Pointer(&values[0])))
}

// More parameters.
const (
	AlLooping      = 0x1007
	AlSampleOffset = 0x1025

	alBits     = 0x2002
	alChannels = 0x2003
	alSize     = 0x2004
)

func (s Source) BuffersProcessed() int32 {
	return s.Geti(AlBuffersProcessed)
}

func (s Source) SetLooping(yes bool) {
	v := int32(0)
	if yes {
		v = 1
	}
	s.Seti(AlLooping, v)
}

func (s Source) GetOffsetSamples() int32 {
	return s.Geti(AlSampleOffset)
}

func (b Buffer) geti(param int32) int32 {
	var v int32
	call(fn.alGetBufferi, uintptr(b), uintptr(param), uintptr(unsafe.Pointer(&v)))
	return v
}

func (b Buffer) GetBits() uint32     { return uint32(b.geti(alBits)) }
func (b Buffer) GetChannels() uint32 { return uint32(b.geti(alChannels)) }
func (b Buffer) GetSize() uint32     { return uint32(b.geti(alSize)) }
