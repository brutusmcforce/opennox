package lobby

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/netip"
	"strings"
)

var _ http.Handler = (*Server)(nil)

// Server is an HTTP Nox lobby server.
type Server struct {
	l         Lobby
	nat       *NATService
	mux       *http.ServeMux
	trustAddr bool // trust IP sent by a remote
	// trustProxy takes the remote IP from the X-Real-IP or X-Forwarded-For headers of requests
	// that come from the local host (a reverse proxy).
	trustProxy bool
}

// SetTrustProxy makes the server use the client address set by a reverse proxy on the same host
// (the X-Real-IP header, or the last address in X-Forwarded-For). Requests from other hosts are not affected,
// so the headers can't be used to register a game for someone else's address.
func (api *Server) SetTrustProxy(v bool) {
	api.trustProxy = v
}

// IPResp represents a response to the address request.
type IPResp struct {
	IP string `json:"ip"`
}

// ServerListResp represents a response to the server list request.
type ServerListResp []GameInfo

// Response wraps all other HTTP responses to separate errors from the rest of the response.
type Response struct {
	Result interface{} `json:"data,omitempty"`
	Err    string      `json:"error,omitempty"`
}

// NewServer creates a new http.Handler from a Lobby implementation.
func NewServer(l Lobby) *Server {
	api := &Server{l: l, mux: http.NewServeMux()}
	api.mux.HandleFunc("/api/v0/address", api.Address)
	api.mux.HandleFunc("/api/v0/games/list", api.ServersList)
	api.mux.HandleFunc("/api/v0/games/register", api.RegisterServer)
	api.mux.HandleFunc("/api/v0/nat", api.NATInfo)
	api.mux.HandleFunc("/api/v0/nat/relay", api.NATRelay)
	return api
}

// SetNAT sets the NAT traversal service, which the NAT endpoints use.
func (api *Server) SetNAT(n *NATService) {
	api.nat = n
}

// NATInfoResp is a response to the NAT info request.
type NATInfoResp struct {
	// UDPPort is the port of the NAT service, on the lobby's host.
	UDPPort uint16 `json:"udp_port"`
	// Relay is set if the lobby can relay games.
	Relay bool `json:"relay,omitempty"`
}

// NATRelayReq is a request for a relay session to a game.
type NATRelayReq struct {
	// Address and Port of the game, as the lobby lists it.
	Address string `json:"addr"`
	Port    int    `json:"port"`
}

// NATRelayResp is a response to the relay request.
type NATRelayResp struct {
	// Port is the UDP port of the relay session, on the lobby's host. The client uses it as the game address.
	Port uint16 `json:"port"`
}

func (api *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cntRequests.WithLabelValues(requestLabels(r)...).Inc()
	api.mux.ServeHTTP(w, r)
}

// jsonResponse writes response, wrapping it into JSON format.
func (api *Server) jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	if code == 0 {
		code = http.StatusOK
	}
	if data == nil {
		// result field in JSON is omitempty, but we want to keep it for this case
		data = (*struct{})(nil)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&Response{
		Result: data,
	})
}

// jsonError writes an error, wrapping it into JSON format.
func (api *Server) jsonError(w http.ResponseWriter, code int, err error) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	if err == nil {
		err = errors.New(http.StatusText(code))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&Response{
		Err: err.Error(),
	})
}

func (api *Server) getAddress(r *http.Request) (string, error) {
	if r.RemoteAddr == "" {
		return "", errors.New("cannot detect IP address")
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !api.trustProxy {
		return ip, err
	}
	if peer, err := netip.ParseAddr(ip); err != nil || !peer.IsLoopback() {
		return ip, nil
	}
	if v, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Real-IP"))); err == nil {
		return v.Unmap().String(), nil
	}
	if h := r.Header.Values("X-Forwarded-For"); len(h) != 0 {
		list := strings.Split(h[len(h)-1], ",")
		if v, err := netip.ParseAddr(strings.TrimSpace(list[len(list)-1])); err == nil {
			return v.Unmap().String(), nil
		}
	}
	return ip, nil
}

func (api *Server) Address(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		ip, err := api.getAddress(r)
		if err != nil {
			api.jsonError(w, http.StatusBadRequest, err)
			return
		}
		api.jsonResponse(w, 0, IPResp{IP: ip})
	default:
		api.jsonError(w, http.StatusMethodNotAllowed, nil)
	}
}

func (api *Server) RegisterServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// set limit to avoid giant requests - 1MB
		body := http.MaxBytesReader(w, r.Body, 1024*1024)
		var req Game
		if err := json.NewDecoder(body).Decode(&req); err != nil {
			api.jsonError(w, http.StatusBadRequest, err)
			return
		}
		if !api.trustAddr {
			addr, err := api.getAddress(r)
			if err != nil {
				api.jsonError(w, http.StatusBadRequest, err)
				return
			}
			req.Address = addr
		}
		req.Map = strings.ToLower(req.Map)
		err := api.l.RegisterGame(r.Context(), &req)
		if err != nil {
			api.jsonError(w, http.StatusBadRequest, err)
			return
		}
		api.jsonResponse(w, 0, nil)
	default:
		api.jsonError(w, http.StatusMethodNotAllowed, nil)
	}
}

func (api *Server) ServersList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		list, err := api.l.ListGames(r.Context())
		if err != nil {
			api.jsonError(w, http.StatusInternalServerError, err)
			return
		}
		api.jsonResponse(w, 0, ServerListResp(list))
	default:
		api.jsonError(w, http.StatusMethodNotAllowed, nil)
	}
}

func (api *Server) NATInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		if api.nat == nil {
			api.jsonError(w, http.StatusNotFound, errors.New("no NAT service"))
			return
		}
		api.jsonResponse(w, 0, NATInfoResp{UDPPort: api.nat.Port(), Relay: api.nat.HasRelay()})
	default:
		api.jsonError(w, http.StatusMethodNotAllowed, nil)
	}
}

func (api *Server) NATRelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.jsonError(w, http.StatusMethodNotAllowed, nil)
		return
	}
	if api.nat == nil {
		api.jsonError(w, http.StatusNotFound, errors.New("no NAT service"))
		return
	}
	var req NATRelayReq
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
		api.jsonError(w, http.StatusBadRequest, err)
		return
	}
	host, err := netip.ParseAddr(req.Address)
	if err != nil || req.Port <= 0 || req.Port > 0xffff {
		api.jsonError(w, http.StatusBadRequest, errors.New("invalid game address"))
		return
	}
	caddr, err := api.getAddress(r)
	if err != nil {
		api.jsonError(w, http.StatusBadRequest, err)
		return
	}
	client, err := netip.ParseAddr(caddr)
	if err != nil {
		api.jsonError(w, http.StatusBadRequest, err)
		return
	}
	port, err := api.nat.OpenRelay(client, netip.AddrPortFrom(host, uint16(req.Port)))
	switch {
	case errors.Is(err, ErrNATRelayLimit):
		api.jsonError(w, http.StatusServiceUnavailable, err)
		return
	case err != nil:
		api.jsonError(w, http.StatusBadRequest, err)
		return
	}
	api.jsonResponse(w, 0, NATRelayResp{Port: port})
}
