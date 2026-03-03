package topology

import (
	"encoding/json"
	"os"
)

type Topology struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	ID string `json:"id"`
}

type Link struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func Parse(path string) (*Topology, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var topology Topology
	if err := json.NewDecoder(f).Decode(&topology); err != nil {
		return nil, err
	}

	return &topology, nil
}
