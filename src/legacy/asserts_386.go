//go:build 386

package legacy

// Layout assertions: these structs keep the layout of the original 32-bit binary.
// On other architectures pointers are larger and the layout differs.

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/internal/protection"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var _ = [1]struct{}{}[32-unsafe.Sizeof(AudioSample{})]
var _ [692 - unsafe.Sizeof(chatBubble{})]byte
var _ [unsafe.Sizeof(chatBubble{}) - 692]byte
var _ [636 - unsafe.Offsetof(chatBubble{}.Length)]byte
var _ [unsafe.Offsetof(chatBubble{}.Length) - 636]byte
var _ [668 - unsafe.Offsetof(chatBubble{}.Drawable)]byte
var _ [unsafe.Offsetof(chatBubble{}.Drawable) - 668]byte
var _ [684 - unsafe.Offsetof(chatBubble{}.Next)]byte
var _ [unsafe.Offsetof(chatBubble{}.Next) - 684]byte
var (
	_ = [1]struct{}{}[200-unsafe.Sizeof(audioEventMetadata{})]
	_ = [1]struct{}{}[16-unsafe.Offsetof(audioEventMetadata{}.Volume)]
	_ = [1]struct{}{}[88-unsafe.Offsetof(audioEventMetadata{}.Ranking)]
	_ = [1]struct{}{}[112-unsafe.Offsetof(audioEventMetadata{}.BucketNode)]
	_ = [1]struct{}{}[128-unsafe.Offsetof(audioEventMetadata{}.Samples)]
	_ = [1]struct{}{}[192-unsafe.Offsetof(audioEventMetadata{}.SampleCount)]
	_ = [1]struct{}{}[576-unsafe.Sizeof(audioEvent{})]
	_ = [1]struct{}{}[168-unsafe.Offsetof(audioEvent{}.Loaded)]
	_ = [1]struct{}{}[176-unsafe.Offsetof(audioEvent{}.Voice)]
	_ = [1]struct{}{}[184-unsafe.Offsetof(audioEvent{}.Timers)]
	_ = [1]struct{}{}[280-unsafe.Offsetof(audioEvent{}.Serial)]
	_ = [1]struct{}{}[288-unsafe.Offsetof(audioEvent{}.Deadline)]
	_ = [1]struct{}{}[296-unsafe.Offsetof(audioEvent{}.Pending)]
	_ = [1]struct{}{}[304-unsafe.Offsetof(audioEvent{}.Order)]
	_ = [1]struct{}{}[440-unsafe.Offsetof(audioEvent{}.ReloadOrder)]
	_ = [1]struct{}{}[568-unsafe.Offsetof(audioEvent{}.ReloadRemaining)]
	_ = [1]struct{}{}[12-unsafe.Sizeof(audioEventHandle{})]
)
var (
	_ = [1]struct{}{}[28-unsafe.Sizeof(audioStreamFormat{})]
	_ = [1]struct{}{}[288-unsafe.Sizeof(audioStreamCatalog{})]
	_ = [1]struct{}{}[268-unsafe.Offsetof(audioStreamCatalog{}.Bag)]
	_ = [1]struct{}{}[280-unsafe.Offsetof(audioStreamCatalog{}.Active)]
	_ = [1]struct{}{}[36-unsafe.Sizeof(audioStreamEntry{})]
	_ = [1]struct{}{}[28-unsafe.Sizeof(audioStreamBuffer{})]
	_ = [1]struct{}{}[24-unsafe.Sizeof(audioStreamChunk{})]
	_ = [1]struct{}{}[28-unsafe.Sizeof(audioStreamCache{})]
	_ = [1]struct{}{}[84-unsafe.Sizeof(audioStreamCacheEntry{})]
	_ = [1]struct{}{}[24-unsafe.Offsetof(audioStreamCacheEntry{}.Buffer)]
	_ = [1]struct{}{}[52-unsafe.Offsetof(audioStreamCacheEntry{}.Cache)]
	_ = [1]struct{}{}[56-unsafe.Offsetof(audioStreamCacheEntry{}.Format)]
	_ = [1]struct{}{}[40-unsafe.Sizeof(audioStreamDescriptor{})]
	_ = [1]struct{}{}[88-unsafe.Sizeof(audioStreamDevice{})]
	_ = [1]struct{}{}[264-unsafe.Sizeof(audioStreamContext{})]
	_ = [1]struct{}{}[88-unsafe.Offsetof(audioStreamContext{}.Timers)]
	_ = [1]struct{}{}[200-unsafe.Offsetof(audioStreamContext{}.Voices)]
	_ = [1]struct{}{}[224-unsafe.Offsetof(audioStreamContext{}.Period)]
	_ = [1]struct{}{}[256-unsafe.Offsetof(audioStreamContext{}.VoiceAPI)]
	_ = [1]struct{}{}[40-unsafe.Sizeof(audioStreamVoiceAPI{})]
	_ = [1]struct{}{}[312-unsafe.Sizeof(audioStreamVoice{})]
	_ = [1]struct{}{}[132-unsafe.Offsetof(audioStreamVoice{}.Context)]
	_ = [1]struct{}{}[176-unsafe.Offsetof(audioStreamVoice{}.Effective)]
	_ = [1]struct{}{}[288-unsafe.Offsetof(audioStreamVoice{}.Buffer)]
	_ = [1]struct{}{}[128-unsafe.Sizeof(audioStreamSystem{})]
)
var (
	_ = [1]struct{}{}[216-unsafe.Offsetof(audioStreamContext{}.Tick)]
	_ = [1]struct{}{}[276-unsafe.Offsetof(audioStreamVoice{}.DataCallback)]
	_ = [1]struct{}{}[280-unsafe.Offsetof(audioStreamVoice{}.LoopCallback)]
	_ = [1]struct{}{}[284-unsafe.Offsetof(audioStreamVoice{}.EndCallback)]
)
var _ = [1]struct{}{}[16-unsafe.Sizeof(spriteAnimationData{})]
var _ = [1]struct{}{}[56-unsafe.Sizeof(spriteConditionalData{})]
var _ [1056 - unsafe.Sizeof(gui.EntryFieldData{})]byte
var _ [unsafe.Sizeof(gui.EntryFieldData{}) - 1056]byte
var _ [4 - unsafe.Sizeof(gui.RadioButtonData{})]byte
var _ [unsafe.Sizeof(gui.RadioButtonData{}) - 4]byte
var _ [16 - unsafe.Sizeof(gui.SliderData{})]byte
var _ [unsafe.Sizeof(gui.SliderData{}) - 16]byte
var _ [140 - unsafe.Sizeof(uiTradeCell{})]byte
var _ [unsafe.Sizeof(uiTradeCell{}) - 140]byte
var _ [136 - unsafe.Offsetof(uiTradeCell{}.Value)]byte
var _ [unsafe.Offsetof(uiTradeCell{}.Value) - 136]byte
var _ [28 - unsafe.Sizeof(collisionHitRecord{})]byte
var _ [unsafe.Sizeof(collisionHitRecord{}) - 28]byte
var _ [16 - unsafe.Offsetof(collisionHitRecord{}.Normal)]byte
var _ [unsafe.Offsetof(collisionHitRecord{}.Normal) - 16]byte
var _ = [1]struct{}{}[16-unsafe.Sizeof(combatAlly{})]
var _ = [1]struct{}{}[8-unsafe.Sizeof(combatFriend{})]
var _ = [1]struct{}{}[80-unsafe.Sizeof(combatFX{})]
var _ = [1]struct{}{}[64-unsafe.Offsetof(combatFX{}.Next)]
var _ = [1]struct{}{}[20-unsafe.Sizeof(combatHealth{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(noxSystemTime{})]
var _ = [1]struct{}{}[332-unsafe.Sizeof(gui.WindowData{})]
var _ = [1]struct{}{}[68-unsafe.Sizeof(gui.Anim{})]
var _ = [1]struct{}{}[100-unsafe.Offsetof(Nox_gui_server_ent_t{}.StatusVal)]
var _ = [1]struct{}{}[168-unsafe.Offsetof(Nox_gui_server_ent_t{}.Field_42)]
var _ = [1]struct{}{}[172-unsafe.Sizeof(Nox_gui_server_ent_t{})]
var _ [256 - int(unsafe.Sizeof(quickbarRecord{}))]byte
var _ [int(unsafe.Sizeof(quickbarRecord{})) - 256]byte
var _ [204 - int(unsafe.Offsetof(quickbarRecord{}.Current))]byte
var _ [int(unsafe.Offsetof(quickbarRecord{}.Current)) - 204]byte
var _ = [1]struct{}{}[52-unsafe.Sizeof(Nox_screenParticle{})]
var _ [32 - int(unsafe.Sizeof(summonRecord{}))]byte
var _ [int(unsafe.Sizeof(summonRecord{})) - 32]byte
var _ [20 - int(unsafe.Offsetof(summonRecord{}.Size))]byte
var _ [int(unsafe.Offsetof(summonRecord{}.Size)) - 20]byte
var _ [24 - int(unsafe.Offsetof(summonRecord{}.Index))]byte
var _ [int(unsafe.Offsetof(summonRecord{}.Index)) - 24]byte
var _ = [1]struct{}{}[524-unsafe.Sizeof(gui.ScrollListBoxItem{})]
var _ = [1]struct{}{}[56-unsafe.Sizeof(gui.ScrollListBoxData{})]
var _ = [1]struct{}{}[1056-unsafe.Sizeof(gui.EntryFieldData{})]
var _ = [1]struct{}{}[60-unsafe.Sizeof(itemRespawnRecord{})]
var _ = [1]struct{}{}[4-unsafe.Offsetof(itemRespawnRecord{}.obj)]
var _ = [1]struct{}{}[20-unsafe.Offsetof(itemRespawnRecord{}.frame)]
var _ = [1]struct{}{}[52-unsafe.Offsetof(itemRespawnRecord{}.next)]
var _ = [1]struct{}{}[56-unsafe.Offsetof(itemRespawnRecord{}.prev)]
var _ = [1]struct{}{}[12-unsafe.Sizeof(legacyListNode{})]
var _ = [1]struct{}{}[8-unsafe.Offsetof(legacyListNode{}.tag)]
var _ = [1]struct{}{}[52-unsafe.Sizeof(playerGroup{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(playerGroupMember{})]
var _ = [1]struct{}{}[12-unsafe.Offsetof(playerGroup{}.name)]
var _ = [1]struct{}{}[32-unsafe.Offsetof(playerGroup{}.id)]
var _ = [1]struct{}{}[40-unsafe.Offsetof(playerGroup{}.members)]
var _ = [1]struct{}{}[12-unsafe.Offsetof(playerGroupMember{}.index)]
var _ = [1]struct{}{}[32-unsafe.Sizeof(mapQuestRow{})]
var _ [140 - unsafe.Sizeof(mapPolygon{})]byte
var _ [unsafe.Sizeof(mapPolygon{}) - 140]byte
var _ [108 - unsafe.Offsetof(mapPolygon{}.Vertices)]byte
var _ [unsafe.Offsetof(mapPolygon{}.Vertices) - 108]byte
var _ [16 - unsafe.Sizeof(mapPolygonVertex{})]byte
var _ [unsafe.Sizeof(mapPolygonVertex{}) - 16]byte
var _ [12 - unsafe.Sizeof(mapPolygonRemap{})]byte
var _ [unsafe.Sizeof(mapPolygonRemap{}) - 12]byte
var _ = [1]struct{}{}[36-unsafe.Sizeof(Nox_map_list_item{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(binfile.MemFile{})]
var _ = [1]struct{}{}[88-unsafe.Sizeof(server.Modifier{})]
var _ = [1]struct{}{}[144-unsafe.Sizeof(server.ModifierEff{})]
var _ = [1]struct{}{}[28-unsafe.Sizeof(server.MissileUpdateData{})]
var _ = [1]struct{}{}[20-unsafe.Sizeof(server.ElevatorUpdateData{})]
var _ = [1]struct{}{}[36-unsafe.Sizeof(server.MoverUpdateData{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(server.MonsterAnim{})]
var _ = [1]struct{}{}[2200-unsafe.Sizeof(server.MonsterUpdateData{})]
var _ = [1]struct{}{}[556-unsafe.Sizeof(server.PlayerUpdateData{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(server.ClassStats{})]
var _ = [1]struct{}{}[148-unsafe.Sizeof(questProgressRecord{})]
var _ [416 - unsafe.Sizeof(reliableMessage{})]byte
var _ [unsafe.Sizeof(reliableMessage{}) - 416]byte
var _ = [1]struct{}{}[60-unsafe.Sizeof(server.TileDef{})]
var _ = [1]struct{}{}[64-unsafe.Sizeof(serverConfigAllowed{})]
var _ = [1]struct{}{}[96-unsafe.Sizeof(serverConfigBlocked{})]
var _ = [1]struct{}{}[64-unsafe.Offsetof(serverConfigBlocked{}.expires)]
var _ = [1]struct{}{}[72-unsafe.Offsetof(serverConfigBlocked{}.address)]
var (
	_ [64 - unsafe.Sizeof(shopSession{})]byte
	_ [unsafe.Sizeof(shopSession{}) - 64]byte
	_ [16 - unsafe.Sizeof(shopItem{})]byte
	_ [unsafe.Sizeof(shopItem{}) - 16]byte
	_ [28 - unsafe.Sizeof(shopStockEntry{})]byte
	_ [unsafe.Sizeof(shopStockEntry{}) - 28]byte
)
var _ = [1]struct{}{}[12-unsafe.Sizeof(spawnPolicyNode{})]
var _ = [1]struct{}{}[148-unsafe.Sizeof(spawnPolicyMonsterListNode{})]
var _ [72 - unsafe.Sizeof(teamUIRow{})]byte
var _ [unsafe.Sizeof(teamUIRow{}) - 72]byte
var _ [60 - unsafe.Offsetof(teamUIRow{}.code)]byte
var _ [unsafe.Offsetof(teamUIRow{}.code) - 60]byte
var _ = [1]struct{}{}[52-unsafe.Sizeof(voteRecord{})]
var _ = [1]struct{}{}[12332-unsafe.Sizeof(server.WallDef{})]
var (
	_ = [1]struct{}{}[516-unsafe.Sizeof(server.Waypoint{})]
)
var _ [44 - unsafe.Sizeof(worldTileCell{})]byte
var _ [unsafe.Sizeof(worldTileCell{}) - 44]byte

// The C manager stores payloads at words 0/1 and links at words 2/3.
var (
	_ [16 - unsafe.Sizeof(protection.Record{})]byte
	_ [unsafe.Sizeof(protection.Record{}) - 16]byte
	_ [4 - unsafe.Alignof(protection.Record{})]byte
	_ [unsafe.Alignof(protection.Record{}) - 4]byte
	_ [0 - unsafe.Offsetof(protection.Record{}.ID)]byte
	_ [4 - unsafe.Offsetof(protection.Record{}.Value)]byte
	_ [unsafe.Offsetof(protection.Record{}.Value) - 4]byte
	_ [8 - unsafe.Offsetof(protection.Record{}.Next)]byte
	_ [unsafe.Offsetof(protection.Record{}.Next) - 8]byte
	_ [12 - unsafe.Offsetof(protection.Record{}.Prev)]byte
	_ [unsafe.Offsetof(protection.Record{}.Prev) - 12]byte
)
