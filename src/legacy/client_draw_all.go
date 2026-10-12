package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/internal/binfile"
)

func init() {
	registerDrawableDrawCallbacks()
	client.RegisterDraw("StaticDraw", drawableDrawKey(drawKey_nox_thing_static_draw), 1, wrapDrawParseGo(spriteParseStatic))
	client.RegisterDraw("StaticRandomDraw", drawableDrawKey(drawKey_nox_thing_static_random_draw), 2, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, b []byte) bool {
		return spriteParseRandom(o, f, b, false)
	}))
	client.RegisterDraw("DoorDraw", drawableDrawKey(drawKey_nox_thing_door_draw), 2, wrapDrawParseGo(objectDoorDrawParse))
	client.RegisterDraw("AnimateDraw", drawableDrawKey(drawKey_nox_thing_animate_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("ConditionalAnimateDraw", drawableDrawKey(drawKey_nox_thing_cond_animate_draw), 4, wrapDrawParseGo(spriteParseConditional))
	client.RegisterDraw("MonsterGeneratorDraw", drawableDrawKey(drawKey_nox_thing_monster_gen_draw), 4, wrapDrawParseGo(spriteParseConditional))
	client.RegisterDraw("AnimateStateDraw", drawableDrawKey(drawKey_nox_thing_animate_state_draw), 8, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, _ []byte) bool { return spriteParseState(o, f) }))
	client.RegisterDraw("SlaveDraw", drawableDrawKey(drawKey_nox_thing_slave_draw), 2, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, b []byte) bool { return spriteParseRandom(o, f, b, true) }))
	client.RegisterDraw("TriggerDraw", drawableDrawKey(drawKey_nox_thing_trigger_draw), 0, nil)
	client.RegisterDraw("PressurePlateDraw", drawableDrawKey(drawKey_nox_thing_pressure_plate_draw), 0, nil)
	client.RegisterDraw("LightningDraw", drawableDrawKey(drawKey_nox_thing_lightning_draw), 0, nil)
	client.RegisterDraw("ChainLightningBoltDraw", drawableDrawKey(drawKey_nox_thing_chain_lightning_bolt_draw), 0, nil)
	client.RegisterDraw("EnergyBoltDraw", drawableDrawKey(drawKey_nox_thing_energy_bolt_draw), 0, nil)
	client.RegisterDraw("GreenBoltDraw", drawableDrawKey(drawKey_nox_thing_green_bolt_draw), 0, nil)
	client.RegisterDraw("PlasmaDraw", drawableDrawKey(drawKey_nox_thing_plasma_draw), 0, nil)
	client.RegisterDraw("YellowSparkDraw", drawableDrawKey(drawKey_nox_thing_yellow_spark_draw), 0, nil)
	client.RegisterDraw("RedSparkDraw", drawableDrawKey(drawKey_nox_thing_red_spark_draw), 0, nil)
	client.RegisterDraw("BlueSparkDraw", drawableDrawKey(drawKey_nox_thing_blue_spark_draw), 0, nil)
	client.RegisterDraw("CyanSparkDraw", drawableDrawKey(drawKey_nox_thing_cyan_spark_draw), 0, nil)
	client.RegisterDraw("GreenSparkDraw", drawableDrawKey(drawKey_nox_thing_green_spark_draw), 0, nil)
	client.RegisterDraw("VioletSparkDraw", drawableDrawKey(drawKey_nox_thing_violet_spark_draw), 0, nil)
	client.RegisterDraw("WhiteSparkDraw", drawableDrawKey(drawKey_nox_thing_white_spark_draw), 0, nil)
	client.RegisterDraw("BlueRainSparkDraw", drawableDrawKey(drawKey_nox_thing_blue_rain_spark_draw), 0, nil)
	client.RegisterDraw("DeathBallSparkDraw", drawableDrawKey(drawKey_nox_thing_death_ball_spark_draw), 0, nil)
	client.RegisterDraw("PixieDraw", drawableDrawKey(drawKey_nox_thing_pixie_draw), 0, nil)
	client.RegisterDraw("PixieDustDraw", drawableDrawKey(drawKey_nox_thing_pixie_dust_draw), 0, nil)
	client.RegisterDraw("ParticleDraw", drawableDrawKey(drawKey_nox_thing_particle_draw), 0, nil)
	client.RegisterDraw("BubbleDraw", drawableDrawKey(drawKey_nox_thing_bubble_draw), 0, nil)
	client.RegisterDraw("VortexDraw", drawableDrawKey(drawKey_nox_thing_vortex_draw), 0, nil)
	client.RegisterDraw("BlackPowderDraw", drawableDrawKey(drawKey_nox_thing_black_powder_draw), 0, nil)
	client.RegisterDraw("SpiderSpitDraw", drawableDrawKey(drawKey_nox_thing_spider_spit_draw), 0, nil)
	client.RegisterDraw("GlyphDraw", drawableDrawKey(drawKey_nox_thing_glyph_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("BoulderDraw", drawableDrawKey(drawKey_nox_thing_boulder_draw), 2, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, b []byte) bool { return spriteParseRandom(o, f, b, true) }))
	client.RegisterDraw("DrainManaDraw", drawableDrawKey(drawKey_nox_thing_drain_mana_draw), 0, nil)
	client.RegisterDraw("GlowOrbDraw", drawableDrawKey(drawKey_nox_thing_glow_orb_draw), 0, nil)
	client.RegisterDraw("GlowOrbMoveDraw", drawableDrawKey(drawKey_nox_thing_glow_orb_move_draw), 0, nil)
	client.RegisterDraw("MagicDraw", drawableDrawKey(drawKey_nox_thing_magic_draw), 0, nil)
	client.RegisterDraw("MagicMissileDraw", drawableDrawKey(drawKey_nox_thing_magic_missle_draw), 0, nil)
	client.RegisterDraw("MagicTailLinkDraw", drawableDrawKey(drawKey_nox_thing_magic_tail_link_draw), 0, nil)
	client.RegisterDraw("MagicMissileTailLinkDraw", drawableDrawKey(drawKey_nox_thing_magic_missle_tail_link_draw), 0, nil)
	client.RegisterDraw("MagicSparkleDraw", drawableDrawKey(drawKey_nox_thing_magic_sparkle_draw), 0, nil)
	client.RegisterDraw("PlayerWaypointDraw", drawableDrawKey(drawKey_nox_thing_player_waypoint_draw), 0, nil)
	client.RegisterDraw("WeaponDraw", drawableDrawKey(drawKey_nox_thing_weapon_draw), 1, wrapDrawParseGo(spriteParseStatic))
	client.RegisterDraw("ArmorDraw", drawableDrawKey(drawKey_nox_thing_armor_draw), 1, wrapDrawParseGo(spriteParseStatic))
	client.RegisterDraw("WeaponAnimateDraw", drawableDrawKey(drawKey_nox_thing_weapon_animate_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("ArmorAnimateDraw", drawableDrawKey(drawKey_nox_thing_armor_animate_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("FlagDraw", drawableDrawKey(drawKey_nox_thing_flag_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("BaseDraw", drawableDrawKey(drawKey_nox_thing_base_draw), 1, wrapDrawParseGo(spriteParseStatic))
	client.RegisterDraw("SphericalShieldDraw", drawableDrawKey(drawKey_nox_thing_spherical_shield_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("SummonEffectDraw", drawableDrawKey(drawKey_nox_thing_summon_effect_draw), 3, wrapDrawParseGo(spriteParseAnimate))
	client.RegisterDraw("UndeadKillerDraw", drawableDrawKey(drawKey_nox_thing_undead_killer_draw), 0, nil)
	client.RegisterDraw("ArrowDraw", drawableDrawKey(drawKey_nox_thing_arrow_draw), 2, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, b []byte) bool { return spriteParseRandom(o, f, b, true) }))
	client.RegisterDraw("WeakArrowDraw", drawableDrawKey(drawKey_nox_thing_weak_arrow_draw), 2, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, b []byte) bool { return spriteParseRandom(o, f, b, true) }))
	client.RegisterDraw("ArrowTailLinkDraw", drawableDrawKey(drawKey_nox_thing_arrow_tail_link_draw), 0, nil)
	client.RegisterDraw("WeakArrowTailLinkDraw", drawableDrawKey(drawKey_nox_thing_weak_arrow_tail_link_draw), 0, nil)
	client.RegisterDraw("BlueRainDraw", drawableDrawKey(drawKey_nox_thing_blue_rain_draw), 0, nil)
	client.RegisterDraw("LevelUpDraw", drawableDrawKey(drawKey_nox_thing_levelup_draw), 0, nil)
	client.RegisterDraw("OblivionUpDraw", drawableDrawKey(drawKey_nox_thing_oblivion_up_draw), 0, nil)
	client.RegisterDraw("RainOrbDraw", drawableDrawKey(drawKey_nox_thing_rain_orb_draw), 0, nil)
	client.RegisterDraw("HarpoonDraw", drawableDrawKey(drawKey_nox_thing_harpoon_draw), 2, wrapDrawParseGo(func(o *client.ObjectType, f *binfile.MemFile, b []byte) bool { return spriteParseRandom(o, f, b, true) }))
	client.RegisterDraw("HarpoonRopeDraw", drawableDrawKey(drawKey_nox_thing_harpoon_rope_draw), 0, nil)
}

var (
	Nox_thing_debug_draw          func(vp *noxrender.Viewport, dr *client.Drawable) int
	Nox_thing_monster_draw        func(vp *noxrender.Viewport, dr *client.Drawable) int
	Nox_thing_vector_animate_draw func(vp *noxrender.Viewport, dr *client.Drawable) int
	Nox_thing_animate_state_draw  func(vp *noxrender.Viewport, dr *client.Drawable) int
	Nox_thing_player_draw         func(vp *noxrender.Viewport, dr *client.Drawable) int
	Nox_thing_npc_draw            func(vp *noxrender.Viewport, dr *client.Drawable) int
)

func Nox_xxx_drawObject_4C4770_draw(vp *noxrender.Viewport, dr *client.Drawable, img noxrender.ImageHandle) {
	objectRenderDraw(vp, dr, img)
}
func Sub_495180(id int) (cur, max int, alt, ok bool) {
	var (
		curV, maxV uint16
		altV       byte
	)
	ok = combatAllyRead(uint32(id), &curV, &maxV, &altV)
	return int(curV), int(maxV), altV != 0, ok
}

func Get_nox_thing_static_random_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_static_random_draw)
}
func Get_nox_thing_red_spark_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_red_spark_draw)
}
func Get_nox_thing_blue_spark_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_blue_spark_draw)
}
func Get_nox_thing_yellow_spark_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_yellow_spark_draw)
}
func Get_nox_thing_green_spark_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_green_spark_draw)
}
func Get_nox_thing_cyan_spark_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_cyan_spark_draw)
}
func Get_nox_thing_monster_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_monster_draw)
}
func Get_nox_thing_maiden_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_maiden_draw)
}
func Get_nox_thing_player_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_player_draw)
}
func Get_nox_thing_animate_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_animate_draw)
}
func Get_nox_thing_vector_animate_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_vector_animate_draw)
}
func Get_nox_thing_released_soul_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_released_soul_draw)
}
func Get_nox_thing_debug_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_debug_draw)
}
func Get_nox_thing_npc_draw() unsafe.Pointer {
	return drawableDrawKey(drawKey_nox_thing_npc_draw)
}
