package command

import (
	"io"
	"os"
)

type Cli interface {
	Out() io.Writer
	Err() io.Writer
	In() io.Reader
	// TODO: client to usernetd (using gRPC)
}

type UsernetCli struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func NewUsernetCli() *UsernetCli {
	return &UsernetCli{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
}

func (c *UsernetCli) Out() io.Writer { return c.out }
func (c *UsernetCli) Err() io.Writer { return c.err }
func (c *UsernetCli) In() io.Reader  { return c.in }
