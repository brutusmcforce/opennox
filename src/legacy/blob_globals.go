package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/legacy/timer"
	"github.com/brutusmcforce/opennox/server"
)

// Variables moved out of the original memory blobs (common/memmap) because they hold pointers,
// which do not fit the 32-bit blob layout on 64-bit platforms. Comments name the old blob location.

// XferRecv holds data received with the transfer protocol; entry 1 is the server's MOTD.
// Was 0x5D4594+826048 (sizes), +826056 (data) and +826064 (ready flags).
var XferRecv [2]struct {
	Size  uint32
	Data  unsafe.Pointer
	Ready uint32
}

// uiInventoryDraggedItem is the inventory item being dragged. Was 0x5D4594+1049848.
var uiInventoryDraggedItem *client.Drawable

// Scoreboard windows. Was 0x5D4594+1090052 (groups), +1090060 (columns) and +1090104 (title).
var (
	ScoreboardGroups  [2]*gui.Window
	scoreboardColumns [5][2]*gui.Window
	scoreboardTitle   *gui.Window
)

// Animations of the confusion and shield effects. Was 0x5D4594+1096456 and +1096460.
var ConfusedBirdiesAnim, SphericalShieldAnim *ImageRef

// DialogVoice is the voice file of the current NPC dialog. Was 0x5D4594+1115312.
var DialogVoice *byte

// AudioGlobalTimers is the timer group shared by audio streams. Was 0x5D4594+1193340.
var AudioGlobalTimers *timer.TimerGroup

var (
	chatBubbleList        *chatBubble    // was 0x5D4594+1197368
	combatFeedDefaultIcon unsafe.Pointer // was 0x5D4594+1203828
	combatFXList          *combatFX      // was 0x5D4594+1203872

	presentationShieldList       [9]*client.Drawable  // was 0x5D4594+1217468
	presentationTransientRayList [96]*client.Drawable // was 0x5D4594+1303540
	presentationRayList          [96]*client.Drawable // was 0x5D4594+1303924

	shopCachedList [32]*shopSession // was 0x5D4594+2386364
)

var (
	characterAcceptButton *gui.Window // was 0x5D4594+1307728

	uiAmountCancelFunc unsafe.Pointer // was 0x5D4594+1319100
	uiAmountAccept     unsafe.Pointer // was 0x5D4594+1319160

	objectRenderLastBelow *client.Drawable // was 0x5D4594+1321516
	objectRenderShinyAnim *ImageRef        // was 0x5D4594+1321524
)

var (
	questSettings  unsafe.Pointer // was 0x5D4594+1556152
	QuestWasClosed bool

	reliablePoolPtr unsafe.Pointer // was 0x5D4594+1565508

	projectileChakramOwner *server.Object // was 0x5D4594+1567840
	projectileChakramObj   *server.Object // was 0x5D4594+1567924

	questEligibilityReplenish  [4]unsafe.Pointer // was 0x5D4594+1568308
	questEligibilitySpecialMod unsafe.Pointer    // was 0x5D4594+1568324

	spellDoorLinkSource    *server.Object // was 0x5D4594+2487716
	monsterShieldThreatObj *server.Object // was 0x5D4594+2487956
	temporaryAntiTarget    *server.Object // was 0x5D4594+2488668
	monsterMeleeTarget     *server.Object // was 0x5D4594+2491564
)

// LocalPlayerDrawable is the drawable of the local player's unit. Was 0x852978+8.
var LocalPlayerDrawable *client.Drawable

// ColorTable points to the 17 colors selectable for list rows and team names. Was 0x85B3FC+132.
var ColorTable [17]*uint32
