package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/internal/binfile"
)

var drawableDrawIdentitySlots [68]byte

const (
	drawKey_nox_thing_lightning_draw              = 0
	drawKey_nox_thing_chain_lightning_bolt_draw   = 1
	drawKey_nox_thing_energy_bolt_draw            = 2
	drawKey_nox_thing_green_bolt_draw             = 3
	drawKey_nox_thing_plasma_draw                 = 4
	drawKey_nox_thing_magic_sparkle_draw          = 5
	drawKey_nox_thing_pixie_draw                  = 6
	drawKey_nox_thing_pixie_dust_draw             = 7
	drawKey_nox_thing_blue_rain_spark_draw        = 8
	drawKey_nox_thing_rain_orb_draw               = 9
	drawKey_nox_thing_red_spark_draw              = 10
	drawKey_nox_thing_blue_spark_draw             = 11
	drawKey_nox_thing_cyan_spark_draw             = 12
	drawKey_nox_thing_green_spark_draw            = 13
	drawKey_nox_thing_yellow_spark_draw           = 14
	drawKey_nox_thing_violet_spark_draw           = 15
	drawKey_nox_thing_death_ball_spark_draw       = 16
	drawKey_nox_thing_white_spark_draw            = 17
	drawKey_nox_thing_particle_draw               = 18
	drawKey_nox_thing_glow_orb_draw               = 19
	drawKey_nox_thing_glow_orb_move_draw          = 20
	drawKey_nox_thing_door_draw                   = 21
	drawKey_nox_thing_arrow_draw                  = 22
	drawKey_nox_thing_weak_arrow_draw             = 23
	drawKey_nox_thing_arrow_tail_link_draw        = 24
	drawKey_nox_thing_weak_arrow_tail_link_draw   = 25
	drawKey_nox_thing_glyph_draw                  = 26
	drawKey_nox_thing_summon_effect_draw          = 27
	drawKey_nox_thing_weapon_draw                 = 28
	drawKey_nox_thing_weapon_animate_draw         = 29
	drawKey_nox_thing_armor_draw                  = 30
	drawKey_nox_thing_armor_animate_draw          = 31
	drawKey_nox_thing_spherical_shield_draw       = 32
	drawKey_nox_thing_monster_gen_draw            = 33
	drawKey_nox_thing_pressure_plate_draw         = 34
	drawKey_nox_thing_trigger_draw                = 35
	drawKey_nox_thing_black_powder_draw           = 36
	drawKey_nox_thing_base_draw                   = 37
	drawKey_nox_thing_flag_draw                   = 38
	drawKey_nox_thing_magic_draw                  = 39
	drawKey_nox_thing_magic_missle_draw           = 40
	drawKey_nox_thing_magic_missle_tail_link_draw = 41
	drawKey_nox_thing_magic_tail_link_draw        = 42
	drawKey_nox_thing_drain_mana_draw             = 43
	drawKey_nox_thing_bubble_draw                 = 44
	drawKey_nox_thing_blue_rain_draw              = 45
	drawKey_nox_thing_levelup_draw                = 46
	drawKey_nox_thing_oblivion_up_draw            = 47
	drawKey_nox_thing_spider_spit_draw            = 48
	drawKey_nox_thing_vortex_draw                 = 49
	drawKey_nox_thing_animate_draw                = 50
	drawKey_nox_thing_cond_animate_draw           = 51
	drawKey_nox_thing_static_draw                 = 52
	drawKey_nox_thing_static_random_draw          = 53
	drawKey_nox_thing_slave_draw                  = 54
	drawKey_nox_thing_boulder_draw                = 55
	drawKey_nox_thing_debug_draw                  = 56
	drawKey_nox_thing_monster_draw                = 57
	drawKey_nox_thing_vector_animate_draw         = 58
	drawKey_nox_thing_released_soul_draw          = 59
	drawKey_nox_thing_animate_state_draw          = 60
	drawKey_nox_thing_player_draw                 = 61
	drawKey_nox_thing_npc_draw                    = 62
	drawKey_nox_thing_harpoon_draw                = 63
	drawKey_nox_thing_harpoon_rope_draw           = 64
	drawKey_nox_thing_undead_killer_draw          = 65
	drawKey_nox_thing_maiden_draw                 = 66
	drawKey_nox_thing_player_waypoint_draw        = 67
)

func drawableDrawKey(id int) unsafe.Pointer {
	return unsafe.Pointer(&drawableDrawIdentitySlots[id])
}

func registerDrawableDrawCallbacks() {
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_lightning_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectLightningDraw(vp, dr, 0))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_chain_lightning_bolt_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectLightningDraw(vp, dr, 0))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_energy_bolt_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectLightningDraw(vp, dr, 1))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_green_bolt_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectLightningDraw(vp, dr, 2))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_plasma_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectPlasmaDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_magic_sparkle_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectMagicSparkle(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_pixie_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectPixie(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_pixie_dust_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectPixieDust(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_blue_rain_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectBlueRainSpark(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_rain_orb_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectRainOrb(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_red_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 0))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_blue_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 1))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_cyan_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 2))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_green_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 3))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_yellow_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 4))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_violet_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 5))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_death_ball_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 6))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_white_spark_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectColoredSpark(vp, dr, 7))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_particle_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		// particle effects are implemented in Go, see Nox_client_particleFxDraw
		if Nox_client_particleFxDraw != nil {
			Nox_client_particleFxDraw(dr)
		}
		return 1
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_glow_orb_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectOrb(vp, dr, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_glow_orb_move_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(effectOrb(vp, dr, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_door_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectDoorDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_arrow_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectArrowDraw(vp, dr, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_weak_arrow_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectArrowDraw(vp, dr, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_arrow_tail_link_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectArrowTailDraw(vp, dr, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_weak_arrow_tail_link_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectArrowTailDraw(vp, dr, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_glyph_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectGlyphDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_summon_effect_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectSummonDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_weapon_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectEquipmentDraw(vp, dr, false, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_weapon_animate_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectEquipmentDraw(vp, dr, false, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_armor_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectEquipmentDraw(vp, dr, true, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_armor_animate_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectEquipmentDraw(vp, dr, true, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_spherical_shield_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectShieldDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_monster_gen_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectGeneratorDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_pressure_plate_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectPressureDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_trigger_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectTriggerDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_black_powder_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectPowderDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_base_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		objectEquipmentDraw(vp, dr, false, false)
		return 1
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_flag_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(objectFlagDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_magic_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleMagicDraw(vp, dr, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_magic_missle_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleMagicDraw(vp, dr, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_magic_missle_tail_link_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleTailDraw(vp, dr, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_magic_tail_link_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleTailDraw(vp, dr, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_drain_mana_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(1)
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_bubble_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleBubbleDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_blue_rain_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleBlueRain(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_levelup_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleLevelUp(vp, dr, false))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_oblivion_up_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleLevelUp(vp, dr, true))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_spider_spit_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleSpiderSpit(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_vortex_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(particleVortexDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_animate_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteAnimateDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_cond_animate_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteConditionalDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_static_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteStaticDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_static_random_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteSlaveDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_slave_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteSlaveDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_boulder_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteBoulderDraw(vp, dr))
	})

	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_debug_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_debug_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_monster_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_monster_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_vector_animate_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_vector_animate_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_released_soul_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_vector_animate_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_animate_state_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_animate_state_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_player_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_player_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_npc_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(Nox_thing_npc_draw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_harpoon_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(spriteSlaveDraw(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_harpoon_rope_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(screenHarpoonRope(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_undead_killer_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(screenUndead(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_maiden_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(screenMaiden(vp, dr))
	})
	client.RegisterDrawableDrawCallbackGo(drawableDrawKey(drawKey_nox_thing_player_waypoint_draw), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return int32(screenWaypoint(vp, dr))
	})
}

func wrapDrawParseGo(fnc func(*client.ObjectType, *binfile.MemFile, []byte) bool) client.ThingFieldFunc {
	return func(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
		StrNCopyBytes(buf, str)
		fnc(typ, f, unsafe.Slice(&buf[0], 256))
		return nil
	}
}

func objectDoorDrawParse(obj *client.ObjectType, f *binfile.MemFile, scratch []byte) bool {
	obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_door_draw)
	obj.DrawData = spriteStaticRandomData(f, scratch)
	return obj.DrawData != nil
}
