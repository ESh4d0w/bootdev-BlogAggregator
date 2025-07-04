package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
	"github.com/google/uuid"
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

	dbUser, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Not found in database: %v", err)
	}

	err = s.cfg.SetUser(dbUser.Name)
	if err != nil {
		return fmt.Errorf("Error setting User: %v\n", err)
	}

	log.Printf("UserName set to %s\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handerRegister doesn't have a state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Register expects one argument: <username>\n")
	}

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("Cannot create User in Database: %v", err)
	}

	err = s.cfg.SetUser(newUser.Name)
	if err != nil {
		return fmt.Errorf("Error setting User: %v\n", err)
	}

	log.Printf("Successfully created User %v", newUser.Name)
	return nil

}

func handlerUserList(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handlerUserList doesn't have a state")
	}
	if len(cmd.args) != 0 {
		return fmt.Errorf("Users expects no arguments\n")
	}

	userList, err := s.db.GetUserList(context.Background())
	if err != nil {
		return fmt.Errorf("Cannot get users from Database: %v", err)
	}
	for _, user := range userList {
		if user.Name == s.cfg.CurrentUserName {
			log.Printf("* %s (current)", user.Name)
			continue
		}
		log.Printf("* %s", user.Name)
	}

	return nil

}

func handlerReset(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handlreReset doesn't have a state")
	}
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset expect 0 arguments")
	}

	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Reset Failed: %v", err)
	}

	log.Printf("Successfully reset")
	return nil

}
