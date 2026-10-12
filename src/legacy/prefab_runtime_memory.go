package legacy

import "unsafe"

// MapPrefabFreeNode releases a prefab cache wrapper or raw payload with the same allocator
// used by its constructors and removal paths. Group references have their
// own tracked Go allocator and must not use this function.
func MapPrefabFreeNode(p unsafe.Pointer) { mapRoomRelease(p) }
