package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFeedFollowsAdd(s *state, cmd command, user database.User) error {
	if s == nil {
		panic("handler_feed_follows.go:handlerFeedFollowsAdd doesn't have state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Expect one argument <url>\n")
	}
	feed, err := s.db.FeedsGetByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Feed not found in database: %v", err)
	}

	feedfollows, err := s.db.FeedFollowsCreate(context.Background(), database.FeedFollowsCreateParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't add feed for user: %v", err)
	}
	log.Printf("%-20s follows %-20s", feedfollows.UserName, feedfollows.FeedName)
	return nil
}

func handlerFeedFollowsForUser(s *state, cmd command, user database.User) error {
	if s == nil {
		panic("handler_feed_follows.go:handlerFeedFollowsList doesn't have state")
	}
	if len(cmd.args) != 0 {
		log.Print("FeedFollowsList doesn't need an argument\n")
	}

	ffuser, err := s.db.FeedFollowsGetForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Recieving feeds for User error: %v", err)
	}

	if len(ffuser) == 0 {
		log.Print("You don't follow any feed!\n")
		return nil
	}

	for _, feed := range ffuser {
		log.Printf("%-20s follows %-20s", feed.UserName, feed.FeedName)
	}
	return nil
}

func handlerFeedFollowsDelete(s *state, cmd command, user database.User) error {
	if s == nil {
		panic("handler_feed_follows.go:handlerFeedFollowsDelete doesn't have state")
	}
	if len(cmd.args) != 1 {
		log.Print("Expect one argument <url>\n")
	}

	feed, err := s.db.FeedsGetByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("You don't follow a feed with that url: %v\n", err)
	}

	err = s.db.FeedFollowsDelete(
		context.Background(),
		database.FeedFollowsDeleteParams{
			FeedID: feed.ID,
			UserID: user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("Couldn't unfollow %v\n", err)
	}
	log.Printf("Success!\nUnfollowed %s", feed.Url)
	return nil

}
