package cli

import "log"

type Run struct{}

var _ Command = (*Run)(nil) // enforcement of an interface

func (s *Run) Name() string {
	return "settings"
}

func (s *Run) Description() string {
	return "display configuration run"
}

func (s *Run) Execute(args []string) error {
	log.Println("----- [Run] -----")
	return nil
}
