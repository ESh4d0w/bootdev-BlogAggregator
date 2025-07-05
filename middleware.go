package main

import (
	"context"
	"fmt"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if s == nil {
			panic("middlewareLoggedIn has no state")
		}

		user, err := s.db.UserGetByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Not found in database: %v", err)
		}

		return handler(s, cmd, user)
	}
}
