package network

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vishvananda/netlink"
	"github.com/white-echidna/usernet/internal/topology"
	"github.com/white-echidna/usernet/internal/network/mocks"
)

func TestNetworkCreationWithMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a topology with two nodes and a link
	topo := &topology.Topology{
		Nodes: []topology.Node{
			{ID: "node1"},
			{ID: "node2"},
		},
		Links: []topology.Link{
			{From: "node1", To: "node2"},
		},
	}

	// Create mock objects
	mockNetlink := mocks.NewMockNetlinkManager(ctrl)

	// Set up expectations
	mockNetlink.EXPECT().LinkAdd(gomock.Any()).Return(nil).Times(1)
	mockNetlink.EXPECT().LinkByName(gomock.Any()).Return(&netlink.Veth{}, nil).Times(2)
	mockNetlink.EXPECT().LinkSetNsPid(gomock.Any(), gomock.Any()).Return(nil).Times(2)

	// Create the network
	pids := map[string]int{
		"node1": 1,
		"node2": 2,
	}
	nm := NewNetworkManager(mockNetlink)
	if err := nm.CreateNetwork(topo, pids); err != nil {
		t.Fatalf("failed to create network: %v", err)
	}
}
