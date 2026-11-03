package app

import (
	"github.com/white-echidna/usernet/control_plane/internal/topology"
)

type App struct {
	topology *topology.Topology
}

func New(topology *topology.Topology) *App {
	return &App{
		topology,
	}
}

func (a *App) Run() error {
	return nil
}
