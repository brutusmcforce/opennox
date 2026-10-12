package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/client/noxrender"
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

// MainMenuImages holds the images of the main menu decoration table (0x587000+168832, 48 bytes per entry);
// they were stored at +4 of each entry.
var MainMenuImages []noxrender.ImageHandle

var (
	scriptBindingStrings      [1024]unsafe.Pointer     // was 0x973F18+26664
	scoreboardClassNames      [3]*uint16               // was 0x5D4594+1084056
	presentationPhonemeImages [8]noxrender.ImageHandle // was 0x5D4594+1096564
)

// List heads.
var (
	serverConfigAllowedList legacyListNode // was 0x5D4594+371364
	serverConfigBlockedList legacyListNode // was 0x5D4594+371500
	playerGroupsList        legacyListNode // was 0x5D4594+599460
	audioEventList          legacyListNode // was 0x5D4594+840612
	teamUIPlayerList        legacyListNode // was 0x5D4594+1045652
	teamUITeamList          legacyListNode // was 0x5D4594+1045668
	clientSequenceList      legacyListNode // was 0x5D4594+1197340
)

var (
	audioEventBuckets       [6][10]legacyListNode    // was 0x5D4594+839892
	audioEventMetadataTable [1023]audioEventMetadata // was 0x5D4594+840628
)

// Potion slots of the bottom bar; each was a 20-byte record followed by a drawable (0x5D4594+1090296, 536 bytes each).
var uiPotionSlots [3]struct {
	Slot     uiPotionSlot
	Drawable client.Drawable
}

// Viewports used to draw items in the UI.
var (
	uiInventoryViewport noxrender.Viewport // was 0x5D4594+1049732
	uiPotionViewport    noxrender.Viewport // was 0x5D4594+1091908
	uiShopViewport      noxrender.Viewport // was 0x5D4594+1098492
	uiAmountViewport    noxrender.Viewport // was 0x5D4594+1319108
	uiTradeViewport     noxrender.Viewport // was 0x5D4594+1320188
)

var (
	uiInventoryLookupResult uiInventoryLookup // was 0x5D4594+1049788
	uiShopCells             [60]uiShopCell    // was 0x5D4594+1098636
	audioStreamSystemData   audioStreamSystem // was 0x5D4594+1193204
	uiTradeCells            [2][4]uiTradeCell // was 0x5D4594+1319284 and +1320308
)

// edgeAssets are the floor edge definitions. Was 0x85B3FC+28644 (60 bytes each).
var edgeAssets [64]server.TileDef

var (
	browserGameIcons  [12]noxrender.ImageHandle // was 0x5D4594+814556
	browserRegionMaps [4]noxrender.ImageHandle  // was 0x5D4594+814900
)

// browserPopupItems are the servers listed in the server browser's popup. Was 0x5D4594+1307316.
var browserPopupItems []unsafe.Pointer
