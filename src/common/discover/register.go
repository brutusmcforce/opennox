package discover

import (
	"context"
	"sync"
	"time"

	"github.com/brutusmcforce/opennox/lobby"
	"golang.org/x/sync/errgroup"

	"github.com/brutusmcforce/opennox/libs/log"
)

const (
	lobbyInterval = lobby.DefaultTimeout / 3
)

var (
	lobbyLog = log.New("lobby")
)

type GameInfo struct {
	lobby.Game
}

// NewServer registers the server with initial game info.
// It will automatically retry the registration until the server is closed.
func NewServer(info GameInfo) *RegServer {
	ctx, stop := context.WithCancel(context.Background())
	s := &RegServer{
		stop:  stop,
		ticks: make(chan time.Time, 1),
		cur:   info,
	}
	s.wg.Go(func() error {
		ticker := time.NewTicker(lobbyInterval)
		defer ticker.Stop()
		var t time.Time
		for {
			select {
			case <-ctx.Done():
				return nil
			case t = <-ticker.C:
			}
			select {
			case <-ctx.Done():
				return nil
			case s.ticks <- t:
			}
		}
	})
	s.wg.Go(func() error {
		return s.doLobby(ctx)
	})
	return s
}

type RegServer struct {
	wg    errgroup.Group
	stop  func()
	ticks chan time.Time

	mu  sync.RWMutex
	cur GameInfo
}

func (s *RegServer) Update(info GameInfo) {
	s.mu.Lock()
	s.cur = info
	s.mu.Unlock()
	select {
	case s.ticks <- time.Now():
	default:
	}
}

func (s *RegServer) Close() error {
	if s.stop == nil {
		return nil
	}
	s.stop()
	s.stop = nil
	return s.wg.Wait()
}

func (s *RegServer) info() GameInfo {
	s.mu.RLock()
	info := s.cur
	s.mu.RUnlock()
	return info
}

type gameRegHost struct {
	s *RegServer
}

func (h gameRegHost) GameInfo(ctx context.Context) (*lobby.Game, error) {
	info := h.s.info()
	return &info.Game, nil
}

func (s *RegServer) doLobby(ctx context.Context) error {
	cli := newLobbyClient()
	for {
		if err := lobby.KeepRegistered(ctx, cli, s.ticks, gameRegHost{s: s}); err != nil {
			lobbyLog.Println(err)
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}
