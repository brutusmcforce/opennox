package legacy

import "github.com/brutusmcforce/opennox/server"

func pingMinimum() uint32 {
	players := &GetServer().S().Players
	result := ^uint32(0)
	found := false
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		if pl.PlayerIndex() == server.HostPlayerIndex || int32(Sub_554240(pl.PlayerIndex())) <= 0 {
			continue
		}
		// The qualifying read and sampled value are separate in the original code.
		value := uint32(Sub_554240(pl.PlayerIndex()))
		if value < result {
			result = value
		}
		found = true
	}
	if !found {
		return 0
	}
	return result
}

func pingAverage() uint32 {
	players := &GetServer().S().Players
	var sum, count int32
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		if pl.PlayerIndex() == server.HostPlayerIndex || int32(Sub_554240(pl.PlayerIndex())) <= 0 {
			continue
		}
		// Preserve the second read, 32-bit wrapping sum and signed division.
		sum += int32(Sub_554240(pl.PlayerIndex()))
		count++
	}
	if count == 0 {
		return 0
	}
	return uint32(sum / count)
}
