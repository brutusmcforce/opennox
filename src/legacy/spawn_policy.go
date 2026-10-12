package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"math"
	"unsafe"
)

func spawnPolicyDistance(a, b *server.Object) float64 {
	return float64(stateDistance(a, b))
}
func spawnPolicyView(p *server.Object) types.Rectf {
	info := p.UpdateDataPlayer().Player
	w := float64(*(*uint16)(unsafe.Add(unsafe.Pointer(info), 10)))
	h := float64(*(*uint16)(unsafe.Add(unsafe.Pointer(info), 12)))
	x, y := float64(p.PosVec.X), float64(p.PosVec.Y)
	return types.Rectf{Min: types.Pointf{X: float32(x - w - 100), Y: float32(y - h - 100)}, Max: types.Pointf{X: float32(w + x + 100), Y: float32(h + y + 100)}}
}
func spawnPolicyContains(r types.Rectf, p types.Pointf) bool {
	return p.X >= r.Min.X && p.X <= r.Max.X && p.Y >= r.Min.Y && p.Y <= r.Max.Y
}
func spawnPolicyJoined(p *server.Object) uint32 {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(p.UpdateDataPlayer().Player), int(layoutServerPlayer.Off(4792)) /* Field4792 */))
}
func spawnPolicyFarCull() {
	core := GetServer().S()
	for n := spawnPolicyHead(); n != nil; {
		next := n.Next
		far := true
		for p := core.Players.FirstUnit(); p != nil; p = core.Players.NextUnit(p) {
			if spawnPolicyJoined(p) == 1 && *(*uint32)(unsafe.Add(p.UpdateData, 312)) == 0 && spawnPolicyDistance(p, n.Object) < 700 {
				far = false
			}
		}
		if far {
			GetServer().DelayedDelete(n.Object)
		}
		n = next
	}
}
func spawnPolicyCandidate(u, player *server.Object) {
	// C numerically converts float views of these object words.
	cl := uint8(math.Float32frombits(uint32(u.ObjClass)))
	flags := uint32(math.Float32frombits(uint32(u.ObjFlags)))
	if cl&2 != 0 && flags&0x20 == 0 && (flags&0x8000 == 0 || monsterIsZombie(u)) && GetServer().S().MapTraceRay(player.PosVec, u.PosVec, 69) {
		(*memmap.PtrUint32(0x5D4594, 2386208))++
	}
}
func spawnPolicyAdmission(_ *server.Object, pos types.Pointf) int {
	core := GetServer().S()
	limit := floatToInt32(float32(core.Balance.Float("MaxOnscreenMonsterCount")))
	for p := core.Players.FirstUnit(); p != nil; p = core.Players.NextUnit(p) {
		if spawnPolicyJoined(p) == 0 {
			continue
		}
		r := spawnPolicyView(p)
		if !spawnPolicyContains(r, pos) {
			continue
		}
		occupied := memmap.PtrUint32(0x5D4594, 2386208)
		*occupied = 0
		core.Map.EachObjInRect(r, func(u *server.Object) bool { spawnPolicyCandidate(u, p); return true })
		if *occupied >= uint32(limit) {
			return 0
		}
	}
	return 1
}
func spawnPolicyMonsterHead() *spawnPolicyMonsterListNode {
	return (*spawnPolicyMonsterListNode)(memmap.WordPtr(dword_5d4594_2386224))
}
func spawnPolicySetMonsterHead(n *spawnPolicyMonsterListNode) {
	dword_5d4594_2386224 = memmap.PtrWord(unsafe.Pointer(n))
}
func spawnPolicyFindMonster(u *server.Object) *spawnPolicyMonsterListNode {
	for n := spawnPolicyMonsterHead(); n != nil; n = n.Next {
		if n.Object == u {
			return n
		}
	}
	return nil
}
func spawnPolicyClearMonsterList() uint32 {
	for n := spawnPolicyMonsterHead(); n != nil; {
		next := n.Next
		spawnPolicyMonsterListClass().FreeObjectFirst(unsafe.Pointer(n))
		dword_5d4594_2386228--
		n = next
	}
	spawnPolicySetMonsterHead(nil)
	return 0
}
func spawnPolicyVisibleCull() uint32 {
	core := GetServer().S()
	limit := floatToInt32(float32(core.Balance.Float("MaxOnscreenMonsterCount")))
	counters := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386232), 32)
	clear(counters)
	over := false
	for p := core.Players.FirstUnit(); p != nil; p = core.Players.NextUnit(p) {
		ind := p.UpdateDataPlayer().Player.PlayerInd
		counters[ind] = 0
		if spawnPolicyJoined(p) == 0 {
			continue
		}
		rect := spawnPolicyView(p)
		for spawn := spawnPolicyHead(); spawn != nil; spawn = spawn.Next {
			u := spawn.Object
			if !spawnPolicyContains(rect, u.PosVec) || !core.MapTraceRay(p.PosVec, u.PosVec, 69) {
				continue
			}
			n := spawnPolicyFindMonster(u)
			if n == nil {
				n = (*spawnPolicyMonsterListNode)(spawnPolicyMonsterListClass().NewObject())
				if n == nil {
					break
				}
				n.Next = spawnPolicyMonsterHead()
				if n.Next != nil {
					n.Next.Prev = n
				}
				spawnPolicySetMonsterHead(n)
				dword_5d4594_2386228++
				n.Object = u
			}
			n.ClassMask |= uint32(1) << ind
			counters[ind]++
			if counters[ind] > uint32(limit) {
				over = true
			}
			n.Distance[ind] = float32(spawnPolicyDistance(u, p))
		}
	}
	if over {
		var nodes []*spawnPolicyMonsterListNode
		for n := spawnPolicyMonsterHead(); n != nil; n = n.Next {
			sum := float64(0)
			cnt := 0
			for _, v := range n.Distance {
				if v != 0 {
					sum += float64(v)
					cnt++
				}
			}
			n.Average = float32(sum / float64(cnt))
			nodes = append(nodes, n)
		}
		// Preserve the C pairwise exchange ordering, including ties. A general
		// unstable sort could choose different equal-distance deletion victims.
		for i := range nodes {
			for j := i + 1; j < len(nodes); j++ {
				if nodes[j].Average > nodes[i].Average {
					nodes[i], nodes[j] = nodes[j], nodes[i]
				}
			}
		}
		for i, n := range nodes {
			n.Prev, n.Next = nil, nil
			if i > 0 {
				n.Prev = nodes[i-1]
			}
			if i+1 < len(nodes) {
				n.Next = nodes[i+1]
			}
		}
		if len(nodes) > 0 {
			spawnPolicySetMonsterHead(nodes[0])
		}
		for {
			largest := uint32(0)
			ind := 0
			for i, n := range counters {
				if n > largest {
					largest = n
					ind = i
				}
			}
			if int32(largest) <= limit {
				break
			}
			for n := spawnPolicyMonsterHead(); n != nil; {
				next := n.Next
				if n.ClassMask&(uint32(1)<<ind) != 0 {
					GetServer().DelayedDelete(n.Object)
					for i := range counters {
						if n.ClassMask&(uint32(1)<<i) != 0 {
							counters[i]--
						}
					}
					if n.Next != nil {
						n.Next.Prev = n.Prev
					}
					if n.Prev != nil {
						n.Prev.Next = n.Next
					} else {
						spawnPolicySetMonsterHead(n.Next)
					}
					spawnPolicyMonsterListClass().FreeObjectFirst(unsafe.Pointer(n))
					dword_5d4594_2386228--
					if int32(counters[ind]) <= limit {
						break
					}
				}
				n = next
			}
		}
	}
	return spawnPolicyClearMonsterList()
}
func spawnPolicyTick() uint32 {
	core := GetServer().S()
	if core.Frame()%(5*core.TickRate()) == 0 {
		spawnPolicyFarCull()
	}
	result := core.Frame() / 15
	if core.Frame()%15 == 0 {
		result = spawnPolicyVisibleCull()
	}
	return result
}
