package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregation(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handlerAggregation doesn't have a state")
	}
	if len(cmd.args) != 0 {
		return fmt.Errorf("Aggegation expect 0 arguments")
	}

	feed, err := s.rssClinet.FetchFeed(context.Background(), `https://www.wagslane.dev/index.xml`)
	if err != nil {
		return fmt.Errorf("FetchFeed failed :%v", err)
	}
	log.Printf("+%v", feed)

	return nil

}

func handlerFeedsAdd(s *state, cmd command, user database.User) error {
	if s == nil {
		panic("handler.go:handlerFeedsAdd doesn't have a state")
	}
	if len(cmd.args) < 2 {
		return fmt.Errorf("FeedAdd requires 2 Arguments: <name> <url>")
	}

	feed, err := s.db.FeedsCreate(context.Background(), database.FeedsCreateParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Can't create feed: %v", err)
	}

	formatFeed(feed, user)
	handlerFeedFollowsAdd(
		s,
		command{
			name: "handlerFeedsAddForward",
			args: []string{cmd.args[1]},
		},
		user,
	)
	return nil
}

func handlerFeedsList(s *state, cmd command) error {
	if s == nil {
		panic("handler.go:handlerFeedsList doesn't have a state")
	}
	if len(cmd.args) != 0 {
		return fmt.Errorf("FeedsList expects no arguments.")
	}

	feedsList, err := s.db.FeedsGetList(context.Background())
	if err != nil {
		return fmt.Errorf("Can't get feeds List from DB")
	}

	for _, feed := range feedsList {
		user, err := s.db.UserGetByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Couldn't get User %v", feed.UserID)
		}
		formatFeed(feed, user)
	}

	return nil
}

func formatFeed(feed database.Feed, user database.User) {
	log.Printf("==================")
	log.Printf("%-20s: %-20s\n", "ID", feed.ID)
	log.Printf("%-20s: %-20s\n", "Created", feed.CreatedAt)
	log.Printf("%-20s: %-20s\n", "Updated", feed.UpdatedAt)
	log.Printf("%-20s: %-20s\n", "Name", feed.Name)
	log.Printf("%-20s: %-20s\n", "Url", feed.Url)
	log.Printf("%-20s: %-20s\n", "User Name", user.Name)
	log.Printf("==================")
}
