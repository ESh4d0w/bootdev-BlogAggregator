package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
)

func handlerPostBrowse(s *state, cmd command, user database.User) error {
	if s == nil {
		panic("handler.go:handlerPostBrowse doesn't hava a state")
	}

	limit := 2
	if len(cmd.args) == 1 {
		if providedLimit, err := strconv.Atoi(cmd.args[0]); err != nil {
			limit = providedLimit
		} else {
			return fmt.Errorf("PostBrowse expedts: <limit>.\n Provided argument was invalid: %w", err)
		}
	}

	posts, err := s.db.PostsGetForUser(context.Background(), database.PostsGetForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Error getting posts for user %s: %v", user.Name, err)
	}

	log.Printf("Found %d posts for User %s:", len(posts), user.Name)
	for _, post := range posts {
		formatPost(post)
	}
	return nil
}

func formatPost(post database.PostsGetForUserRow) {
	log.Printf("%-10s from %s\n", post.PublishedAt.Time.Format("Mon 02 Jan"), post.FeedName)
	log.Printf("===== %s =====\n", post.Title)
	log.Printf("%v\n", post.Description.String)
	log.Printf("=============\n")
	log.Print("Link: %s", post.Url)
	log.Printf("=============\n")
}
