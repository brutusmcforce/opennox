package server

// MonsterSoundSet is a sound set of monsters (monster.snd).
type MonsterSoundSet struct {
	Name  *byte            // 0, 0
	Slots [18]uint32       // 1, 4: sound IDs by slot (see Sound)
	Next  *MonsterSoundSet // 19, 76
	_     uint32           // 20, 80
}

// Sound returns the sound ID in slot i (1-18, the 386 offset of the slot divided by 4).
func (s *MonsterSoundSet) Sound(i int) uint32 {
	return s.Slots[i-1]
}

// SoundPtr returns a pointer to slot i (see Sound).
func (s *MonsterSoundSet) SoundPtr(i int) *uint32 {
	return &s.Slots[i-1]
}
