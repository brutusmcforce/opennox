package opennox

import (
	"errors"
	"path"
	"strings"

	"github.com/brutusmcforce/opennox/libs/common"
	ns3vm "github.com/brutusmcforce/opennox/noxscript/ns/v3/vm"
	ns4 "github.com/brutusmcforce/opennox/noxscript/ns/v4"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

var (
	_ ns4.Game   = noxScriptImpl{}
	_ ns3vm.Game = noxScriptImpl{}
)

func (s *Server) noxScriptP() noxScriptNS {
	return noxScriptNS{s, s.Server.NoxScriptNS()}
}

func (s *Server) NoxScript() ns4.Implementation {
	return s.noxScriptP()
}

func (s *Server) NoxScriptVM() ns3vm.Implementation {
	return s.NoxScriptNSVM()
}

func (s noxScriptImpl) NoxScript() ns4.Implementation {
	return s.s.NoxScript()
}

func (s noxScriptImpl) NoxScriptVM() ns3vm.Implementation {
	return s.s.NoxScriptVM()
}

type noxScriptNS struct {
	s *Server
	server.NoxScriptNS
}

func (s noxScriptNS) AutoSave() {
	if noxflags.HasGame(noxflags.GameModeCoop) {
		SaveCoop(common.SaveAuto)
	}
}

func (s noxScriptNS) loadMap(name string, opts *ns4.LoadMapOptions) error {
	if opts == nil {
		opts = new(ns4.LoadMapOptions)
	}
	if ext := path.Ext(name); ext != "" {
		return errors.New("map name must not contain file extension")
	}
	if strings.ContainsAny(name, "/\\.") {
		return errors.New("map name must not contain path elements")
	}
	name = path.Clean(name)
	// TODO: support other map loading options
	return noxLoadMap(name, mapLoadOptions{
		Force: opts.IgnoreMapType,
	})
}

func (s noxScriptNS) LoadMap(name string, opts *ns4.LoadMapOptions) {
	if err := s.loadMap(name, opts); err != nil {
		panic(err)
	}
}
