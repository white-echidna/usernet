package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/white-echidna/usernet/internal/network"
	"github.com/white-echidna/usernet/internal/topology"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <topology_file>", os.Args[0])
	}
	topo, err := topology.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("failed to parse topology: %v", err)
	}

	var wg sync.WaitGroup
	pids := make(map[string]int)
	cmds := make(map[string]*exec.Cmd)

	for _, node := range topo.Nodes {
		cmd := exec.Command("/bin/true")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWNET,
		}
		if err := cmd.Start(); err != nil {
			log.Printf("failed to start data plane process for node %s: %v", node.ID, err)
			continue
		}
		pids[node.ID] = cmd.Process.Pid
		cmds[node.ID] = cmd
	}

	nm := network.NewRealNetworkManager()
	if err := nm.CreateNetwork(topo, pids); err != nil {
		log.Fatalf("failed to create network: %v", err)
	}

	for id, cmd := range cmds {
		wg.Add(1)
		go func(id string, cmd *exec.Cmd) {
			defer wg.Done()
			if err := cmd.Wait(); err != nil {
				log.Printf("data plane process for node %s exited with error: %v", id, err)
			}
		}(id, cmd)
	}

	wg.Wait()
}
