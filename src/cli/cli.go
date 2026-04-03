package cli

import (
	"cloaq/src/monitor"
)

type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

var Commands = []Command{
	&Settings{},
	&Run{},
	&Help{},
	&monitor.Monitor{},
}
