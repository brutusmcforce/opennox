package opennox

const (
	NOX_CLIENT_VERS_CODE = noxProtoVersionHighRes
	noxMaxWidth          = 3840
	noxMaxHeight         = 2160
)

// The largest part of the world that can be drawn at once, e.g. when zoomed out (see video_zoom.go).
// It's limited by the size of the light grid.
const (
	worldViewMaxWidth  = (lightGridW - 2) * lightGrid
	worldViewMaxHeight = (lightGridH - 2) * lightGrid
)
