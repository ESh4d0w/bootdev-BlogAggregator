package main

import (
	"context"
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

func handlerReset(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handlreReset doesn't have a state")
	}
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset expect 0 arguments")
	}

	err := s.db.UserReset(context.Background())
	if err != nil {
		return fmt.Errorf("Reset User Failed: %v", err)
	}
	err = s.db.FeedsReset(context.Background())
	if err != nil {
		return fmt.Errorf("Reset Feed Failed: %v", err)
	}
	err = s.db.FeedFollowsReset(context.Background())
	if err != nil {
		return fmt.Errorf("Reset FeedFollows Failed: %v", err)
	}

	log.Printf("Successfully reset")
	return nil

}
