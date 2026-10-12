package lobby

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	sourceOpenNox = "opennox"
	sourceXWIS    = "xwis"
)

var (
	serverLabelNames = []string{"src", "addr", "port", "name", "vers", "mode", "map"}
	cntGameSeen      = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "nox_game_seen",
		Help: "Number of times the game was seen online",
	}, serverLabelNames)
	cntGameExpired = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "nox_game_expired",
		Help: "Number of times the game registration expired",
	}, serverLabelNames)
	cntGamePlayers = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nox_game_players",
		Help: "Number of players in the game",
	}, serverLabelNames)
	cntXWISRooms = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nox_xwis_rooms",
		Help: "Number of XWIS rooms",
	})
	cntXWISGames = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nox_xwis_games",
		Help: "Number of XWIS rooms",
	})
	cntRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "nox_http_requests",
		Help: "Number of HTTP requests to the API",
	}, []string{"method", "endpoint", "agent"})
)

func serverLabels(src string, g *Game) []string {
	return []string{
		src,
		g.Address,
		strconv.Itoa(g.Port),
		g.Name,
		g.Vers,
		string(g.Mode),
		g.Map,
	}
}

// requestLabels returns the metric labels for an API request. All values come from a fixed set:
// the path, method and agent are sent by the client, and every distinct value would be a new time series.
func requestLabels(r *http.Request) []string {
	method := r.Method
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodHead:
	default:
		method = "other"
	}
	endpoint := r.URL.Path
	switch endpoint {
	case "/api/v0/address", "/api/v0/games/list", "/api/v0/games/register":
	default:
		endpoint = "other"
	}
	agent := "other"
	if strings.HasPrefix(r.Header.Get("User-Agent"), "OpenNox/") {
		agent = "opennox"
	}
	return []string{method, endpoint, agent}
}
