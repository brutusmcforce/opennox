package legacy

import (
	"github.com/brutusmcforce/opennox/server"
)

func resetNetworkAliases(table *[255]server.PlayerNetData) {
	*table = [255]server.PlayerNetData{}
}

func selectNetworkAlias(table *[255]server.PlayerNetData, key1, key2 int32, frame uint32) byte {
	start := int(byte(key1))
	if start == 0 || start == 255 {
		start = 1
	}
	slot := start
	for {
		rec := &table[slot]
		// Stored keys are zero-extended for comparison with the original full ints.
		if (int32(rec.Field0) == key1 && int32(rec.Field2) == key2) || rec.Frame4 < frame {
			return byte(slot)
		}
		slot++
		if slot == 255 {
			slot = 1
		}
		if slot == start {
			return 255
		}
	}
}
