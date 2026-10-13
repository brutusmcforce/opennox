package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
)

// Matching allocator boundary for spriteDataAlloc's raw libc calloc allocations.
func spriteDataFree(p unsafe.Pointer) {
	if p != nil {
		legacyFree(p)
	}
}

func spriteVectorFree(frames *[9]*noxrender.ImageHandle) {
	for direction := 0; direction < 8; direction++ {
		slot := direction
		if slot >= 4 {
			slot++
		} // middle word is metadata, not an owned frame list
		spriteDataFree(unsafe.Pointer(frames[slot]))
	}
}

// spriteDataFreeKind frees draw data of the given kind (see client.RegisterDraw).
func spriteDataFreeKind(p unsafe.Pointer, kind int) {
	switch kind {
	case 2, 3: // a frame list (SpriteStaticRandomData, SpriteAnimationData, ...)
		spriteDataFree(unsafe.Pointer(client.SpriteFrames(p)))
	case 4:
		for _, frames := range (*spriteConditionalData)(p).Frames {
			spriteDataFree(unsafe.Pointer(frames))
		}
	case 5:
		spriteVectorFree(&(*client.AnimationVector)(p).Frames)
	case 6:
		dd := (*client.PlayerDrawData)(p)
		for state := range dd.Anim {
			anim := &dd.Anim[state]
			// Naked, Armor and Weapon are one list of equipment animations
			groups := unsafe.Slice(&anim.Naked, 1+len(anim.Armor)+len(anim.Weapon))
			for _, group := range groups {
				if group != nil {
					spriteVectorFree(&group.Frames)
					spriteDataFree(unsafe.Pointer(group))
				}
			}
		}
	case 7:
		dd := (*client.MonsterDrawData)(p)
		for i := range dd.Anim {
			spriteVectorFree(&dd.Anim[i].Frames)
		}
	case 8:
		dd := (*client.SpriteStateData)(p)
		for i := range dd.States {
			spriteVectorFree(&dd.States[i].Frames)
		}
	}
	spriteDataFree(p)
}
