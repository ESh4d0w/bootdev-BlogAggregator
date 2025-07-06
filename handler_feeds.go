package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregation(s *state, cmd command) error {
	if s == nil {
		panic("hander.go:handlerAggregation doesn't have a state")
	}
	if len(cmd.args) != 1 {
		return fmt.Errorf("Aggegation expect 1 arg: <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Invalid <time_between_reqs>: %v", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feedsScrapeFeeds(s)
	}
}

func feedsScrapeFeeds(s *state) {
	if s == nil {
		panic("hander.go:feedsScrapeFeeds doesn't have a state")
	}
	feed, err := s.db.FeedsGetNextToFetch(context.Background())
	if err != nil {
		log.Printf("FeedsGetNext toFetch failed: %v", err)
		return
	}
	log.Println("Found Feed to Fetch")
	feedsScrapeFeed(s, feed)
}

func feedsScrapeFeed(s *state, feed database.Feed) {
	_, err := s.db.FeedsMarkedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark Feed %s as Fetched: %v\n", feed.Name, err)
		return
	}
	rssFeed, err := s.rssClinet.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch Feed %s: %v\n", feed.Name, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid:  true,
		}
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		log.Printf("Found Title: %-20s\n", item.Title)
		_, err := s.db.PostsCreate(context.Background(), database.PostsCreateParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}

	}
	log.Printf("Fetched Feed %s found %d posts\n", feed.Name, len(rssFeed.Channel.Item))
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
	log.Printf("%-20s: %-20v\n", "Last Fetched At", feed.LastFetchedAt)
	log.Printf("==================")
}
