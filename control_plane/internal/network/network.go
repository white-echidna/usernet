package network

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/white-echidna/usernet/internal/topology"
)

type NetlinkManager interface {
	LinkAdd(link netlink.Link) error
	LinkByName(name string) (netlink.Link, error)
	LinkSetNsPid(link netlink.Link, pid int) error
}

type netlinkManager struct{}

func (nl *netlinkManager) LinkAdd(link netlink.Link) error {
	return netlink.LinkAdd(link)
}

func (nl *netlinkManager) LinkByName(name string) (netlink.Link, error) {
	return netlink.LinkByName(name)
}

func (nl *netlinkManager) LinkSetNsPid(link netlink.Link, pid int) error {
	return netlink.LinkSetNsPid(link, pid)
}

type NetworkManager struct {
	netlink NetlinkManager
}

func NewNetworkManager(netlink NetlinkManager) *NetworkManager {
	return &NetworkManager{netlink: netlink}
}

func NewRealNetworkManager() *NetworkManager {
	return NewNetworkManager(&netlinkManager{})
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
	la := netlink.NewLinkAttrs()
	la.Name = veth0
	veth := &netlink.Veth{
		LinkAttrs: la,
		PeerName:  veth1,
	}
	if err := nm.netlink.LinkAdd(veth); err != nil {
		return fmt.Errorf("failed to create veth pair: %w", err)
	}

	link1, err := nm.netlink.LinkByName(veth0)
	if err != nil {
		return fmt.Errorf("failed to get link %s: %w", veth0, err)
	}
	if err := nm.netlink.LinkSetNsPid(link1, pid1); err != nil {
		return fmt.Errorf("failed to move %s to pid %d: %w", veth0, pid1, err)
	}

	link2, err := nm.netlink.LinkByName(veth1)
	if err != nil {
		return fmt.Errorf("failed to get link %s: %w", veth1, err)
	}
	if err := nm.netlink.LinkSetNsPid(link2, pid2); err != nil {
		return fmt.Errorf("failed to move %s to pid %d: %w", veth1, pid2, err)
	}

	return nil
}
