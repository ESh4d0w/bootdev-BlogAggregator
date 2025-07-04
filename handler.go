package main

import (
	"fmt"
	"log"
)

type command struct {
	name string
	args []string
}

type commands struct {
	nameToFunction map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:run doesn't have a state")
	}
	f, ok := c.nameToFunction[cmd.name]
	if !ok {
		return fmt.Errorf("No Function found for name: %s\n", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if name == "" {
		panic("handler.go:register don't have a name")
	}
	c.nameToFunction[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handerLogin doesn't have a state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Login expects one argument: <username>\n")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error setting User: %v\n", err)
	}
	log.Printf("UserName set to %s\n", s.cfg.CurrentUserName)
	return nil
}
