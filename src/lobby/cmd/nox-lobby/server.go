package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/brutusmcforce/opennox/lobby"
)

func init() {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "run the lobby web server",
	}
	fHost := cmd.Flags().String("host", ":8080", "host the server will listen on")
	fMonitor := cmd.Flags().String("monitor", "127.0.0.1:6060", "host the monitoring api will listen on")
	fPprof := cmd.Flags().Bool("pprof", false, "serve pprof debug handlers on the monitoring api")
	fNAT := cmd.Flags().String("nat", "", "UDP address for the NAT traversal service (hole punching), e.g. :18600; empty disables it")
	fRelay := cmd.Flags().String("relay-ports", "", "UDP port range for relay sessions, e.g. 18601-18640; empty disables the relay")
	fRelayPerIP := cmd.Flags().Int("relay-per-ip", 4, "relay sessions per client address")
	fTrustProxy := cmd.Flags().Bool("trust-proxy", false, "take client addresses from X-Real-IP/X-Forwarded-For of requests from localhost (behind a reverse proxy)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var lb lobby.Lobby = lobby.NewLobby()
		lsrv := lobby.NewServer(lb)
		lsrv.SetTrustProxy(*fTrustProxy)
		if *fNAT != "" {
			cfg := lobby.NATConfig{MaxRelaysPerIP: *fRelayPerIP}
			if *fRelay != "" {
				a, b, ok := strings.Cut(*fRelay, "-")
				first, err1 := strconv.ParseUint(a, 10, 16)
				last, err2 := strconv.ParseUint(b, 10, 16)
				if !ok || err1 != nil || err2 != nil || first == 0 || last < first {
					return fmt.Errorf("invalid relay port range: %q", *fRelay)
				}
				cfg.RelayFirst, cfg.RelayLast = uint16(first), uint16(last)
			}
			nat, err := lobby.ListenNAT(*fNAT, cfg)
			if err != nil {
				return err
			}
			defer func() { _ = nat.Close() }()
			lsrv.SetNAT(nat)
			log.Println("serving NAT traversal on udp", *fNAT, "relay ports:", *fRelay)
			go func() {
				if err := nat.Serve(); err != nil {
					log.Println("NAT service:", err)
				}
			}()
		}
		// TODO: auto TLS with Let's Encrypt
		srv := &http.Server{
			Addr:              *fHost,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       2 * time.Minute,
			MaxHeaderBytes:    16 << 10,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println(r.RemoteAddr, r.Method, r.URL)
				lsrv.ServeHTTP(w, r)
			}),
		}
		log.Println("serving lobby on", srv.Addr)
		if *fMonitor != "" {
			mon := http.NewServeMux()
			mon.Handle("/metrics", promhttp.Handler())
			if *fPprof {
				mon.HandleFunc("/debug/pprof/", pprof.Index)
				mon.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
				mon.HandleFunc("/debug/pprof/profile", pprof.Profile)
				mon.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
				mon.HandleFunc("/debug/pprof/trace", pprof.Trace)
			}
			log.Println("serving monitoring on", *fMonitor)
			go func() {
				msrv := &http.Server{Addr: *fMonitor, Handler: mon, ReadHeaderTimeout: 10 * time.Second}
				if err := msrv.ListenAndServe(); err != nil {
					log.Println(err)
				}
			}()
		}
		return srv.ListenAndServe()
	}
	Root.AddCommand(cmd)
}
