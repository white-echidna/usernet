package command

import (
	"io"
	"os"
)

type Cli interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
	// TODO: client to usernetd (using gRPC)
}

type usernetCli struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func NewUsernetCli() *usernetCli {
	return &usernetCli{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
}

func (c *usernetCli) In() io.Reader  { return c.in }
func (c *usernetCli) Out() io.Writer { return c.out }
func (c *usernetCli) Err() io.Writer { return c.err }
