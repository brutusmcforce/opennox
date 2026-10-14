package main

import (
	"log"
	"net/http"
	"net/http/pprof"
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
	fTrustProxy := cmd.Flags().Bool("trust-proxy", false, "take client addresses from X-Real-IP/X-Forwarded-For of requests from localhost (behind a reverse proxy)")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var lb lobby.Lobby = lobby.NewLobby()
		lsrv := lobby.NewServer(lb)
		lsrv.SetTrustProxy(*fTrustProxy)
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
