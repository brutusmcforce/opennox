package ns

// PlayerAction is what a player is doing, see Obj.PlayerAction.
type PlayerAction string

const (
	PlayerIdle      = PlayerAction("IDLE")
	PlayerWalk      = PlayerAction("WALK")
	PlayerRun       = PlayerAction("RUN")
	PlayerAttack    = PlayerAction("ATTACK")
	PlayerCast      = PlayerAction("CAST")
	PlayerJump      = PlayerAction("JUMP")
	PlayerBlock     = PlayerAction("BLOCK")
	PlayerDying     = PlayerAction("DYING")
	PlayerDead      = PlayerAction("DEAD")
	PlayerShakeFist = PlayerAction("SHAKE_FIST")
	PlayerLaugh     = PlayerAction("LAUGH")
	PlayerPoint     = PlayerAction("POINT")
	// PlayerOther is any other action, e.g. picking up an item.
	PlayerOther = PlayerAction("OTHER")
)

// MonsterAction is the current AI action of a monster or NPC, e.g. "ACTION_IDLE", "ACTION_HUNT"
// or "ACTION_MELEE_ATTACK"; see Obj.MonsterAction.
type MonsterAction string
