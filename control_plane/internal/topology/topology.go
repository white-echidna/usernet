package topology

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Topology struct {
	Nodes []Node `yaml:"nodes"`
	Links []Link `yaml:"links"`
}

type Node struct {
	ID string `yaml:"id"`
}

type Link struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

func Parse(path string) (*Topology, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var topology Topology

	if err := yaml.NewDecoder(f).Decode(&topology); err != nil {
		return nil, err
	}

	return &topology, nil
}
