package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/brutusmcforce/opennox/xwis"
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
	fGlobal := cmd.Flags().Bool("global", false, "run the server in a global mode (monitor other servers)")
	fXWIS := cmd.Flags().Bool("xwis", true, "list games from XWIS as well")
	fXLogin := cmd.Flags().String("xlogin", "", "XWIS login to use")
	fXPass := cmd.Flags().String("xpass", "", "XWIS password to use")
	fXCache := cmd.Flags().Duration("xcache", lobby.DefaultTimeout/2, "XWIS cache duration")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var lb lobby.Lobby = lobby.NewLobby()
		if *fXWIS {
			log.Println("logging in to XWIS")
			c, err := xwis.NewClient(context.Background(), *fXLogin, *fXPass)
			if err != nil {
				return err
			}
			defer c.Close()
			var lx lobby.Lister = lobby.NewXWISWithClient(c)
			if *fXCache > 0 {
				lx = lobby.Cache(lx, *fXCache)
			}
			lb = lobby.Overlay(lb, lx)
		}
		lsrv := lobby.NewServer(lb)
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
			if *fGlobal {
				// For the global server, we want to monitor the list of all games.
				// Since XWIS list will only be retrieved when request comes in, we must periodically do list requests here.
				go func() {
					ctx := context.Background()
					ticker := time.NewTicker(lobby.DefaultTimeout / 2)
					defer ticker.Stop()
					for range ticker.C {
						_, _ = lb.ListGames(ctx)
					}
				}()
			}
		}
		return srv.ListenAndServe()
	}
	Root.AddCommand(cmd)
}
