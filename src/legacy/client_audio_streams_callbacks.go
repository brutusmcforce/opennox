package legacy

import (
	"unsafe"
)

// Non-executable identities for the internal stream callbacks.
var audioStreamCallbacks = make(map[unsafe.Pointer]func(unsafe.Pointer) int32)

var audioStreamCallbackKeys [4]byte
var audioBridgeCallbackKeys [16]byte

func audioBridgeCallbackKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&audioBridgeCallbackKeys[id])
}

func audioStreamKnownCallback(fn, arg unsafe.Pointer) (int32, bool) {
	switch fn {
	case unsafe.Pointer(&audioStreamCallbackKeys[0]):
		return audioStreamContextTick((*audioStreamContext)(arg)), true
	case unsafe.Pointer(&audioStreamCallbackKeys[1]):
		return audioStreamVoiceData((*audioStreamVoice)(arg)), true
	case unsafe.Pointer(&audioStreamCallbackKeys[2]):
		return audioStreamVoiceLoop((*audioStreamVoice)(arg)), true
	case unsafe.Pointer(&audioStreamCallbackKeys[3]):
		return audioStreamVoiceEnd((*audioStreamVoice)(arg)), true
	}
	switch fn {
	case audioBridgeCallbackKey(0):
		return int32(Sub_43EC30(arg)), true
	case audioBridgeCallbackKey(1):
		return int32(Sub_43ECB0(arg)), true
	case audioBridgeCallbackKey(2):
		return int32(Sub_43ED00(arg)), true
	case audioBridgeCallbackKey(3):
		return int32(Sub_43EFD0(arg)), true
	case audioBridgeCallbackKey(4):
		p := *(**AudioSample)(unsafe.Add(arg, 272))
		p.Smp.Stop()
		return 0, true
	case audioBridgeCallbackKey(5):
		panic("abort")
	case audioBridgeCallbackKey(6):
		return 0, true
	case audioBridgeCallbackKey(7):
		return int32(Sub_43F060(arg)), true
	case audioBridgeCallbackKey(8):
		return 0, true
	case audioBridgeCallbackKey(9):
		return int32(Sub_43E940(arg)), true
	case audioBridgeCallbackKey(10):
		Sub_43E9F0()
		return 0, true
	case audioBridgeCallbackKey(11):
		return int32(Sub_43EA20(arg)), true
	case audioBridgeCallbackKey(12):
		return int32(Sub_43EC10()), true
	case audioBridgeCallbackKey(13):
		return audioEventVoiceLoop((*audioStreamVoice)(arg)), true
	case audioBridgeCallbackKey(14):
		return audioEventVoiceEnded((*audioStreamVoice)(arg)), true
	case audioBridgeCallbackKey(15):
		return audioEventVoiceStopped((*audioStreamVoice)(arg)), true
	default:
		if callback := audioStreamCallbacks[fn]; callback != nil {
			return callback(arg), true
		}
		return 0, false
	}
}

func AudioStreamCallbackVoid(fn, arg unsafe.Pointer) {
	if _, ok := audioStreamKnownCallback(fn, arg); !ok {
		panic("unregistered audio stream callback")
	}
}

func AudioStreamCallbackInt(fn, arg unsafe.Pointer) int32 {
	if result, ok := audioStreamKnownCallback(fn, arg); ok {
		return result
	}
	panic("unregistered audio stream callback")
}
