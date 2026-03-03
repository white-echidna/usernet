package network

import (
	"github.com/white-echidna/usernet/internal/topology"
)

type NetworkManager struct{}

func NewNetworkManager() *NetworkManager {
	return &NetworkManager{}
}

func (nm *NetworkManager) CreateNetwork(topo *topology.Topology, pids map[string]int) error {
	for _, link := range topo.Links {
		pid1 := pids[link.From]
		pid2 := pids[link.To]
		veth0 := link.From + "-" + link.To
		veth1 := link.To + "-" + link.From

		if err := nm.createVethPair(veth0, veth1, pid1, pid2); err != nil {
			return err
		}
	}
	return nil
}

func (nm *NetworkManager) createVethPair(veth0, veth1 string, pid1, pid2 int) error {
	// TODO: implement raw RTM_NEWLINK payload serialization here.
	return nil
}
