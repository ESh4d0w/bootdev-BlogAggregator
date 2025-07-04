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

func handlerFeedAdd(s *state, cmd command) error {
	if s == nil {
		panic("handler.go:handlerFeedAdd doesn't have a state")
	}
	user, err := s.db.UserGetByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Not found in database: %v", err)
	}
	if len(cmd.args) < 2 {
		return fmt.Errorf("FeedAdd requires 2 Arguments: <name> <url>")
	}

	feed, err := s.db.FeedCreate(context.Background(), database.FeedCreateParams{
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
	formatFeed(feed)
	return nil
}

func handlerFeedList(s *state, cmd command) error {
	if s == nil {
		panic("handler.go:handlerFeedList doesn't have a state")
	}
	if len(cmd.args) != 0 {
		return fmt.Errorf("FeedList expects no arguments.")
	}

	feedList, err := s.db.FeedGetList(context.Background())
	if err != nil {
		return fmt.Errorf("Can't get feed List from DB")
	}

	formatFeedList(feedList)
	return nil
}

func formatFeed(feed database.Feed) {
	log.Printf("%-20s: %-20s\n", "ID", feed.ID)
	log.Printf("%-20s: %-20s\n", "Created", feed.CreatedAt)
	log.Printf("%-20s: %-20s\n", "Updated", feed.UpdatedAt)
	log.Printf("%-20s: %-20s\n", "Name", feed.Name)
	log.Printf("%-20s: %-20s\n", "Url", feed.Url)
	log.Printf("%-20s: %-20s\n", "User ID", feed.UserID)
}

func formatFeedList(feedList []database.FeedGetListRow) {
	for _, feed := range feedList {
		log.Printf("%-15s: %-20s %-20s", feed.UserName.String, feed.Name, feed.Url)
	}
}
