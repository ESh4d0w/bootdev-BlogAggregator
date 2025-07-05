package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/config"
	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
	"github.com/esh4d0w/bootdev-BlogAggregator/internal/rss"
	_ "github.com/lib/pq"
)

type state struct {
	db        *database.Queries
	cfg       *config.Config
	rssClinet *rss.Client
}

func main() {
	log.Printf("Blog Aggregator")

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config Read Error: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	dbQueries := database.New(db)

	rssClient := rss.NewClient(5 * time.Second)

	programState := state{
		db:        dbQueries,
		cfg:       &cfg,
		rssClinet: &rssClient,
	}

	cmds := commands{
		nameToFunction: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUserList)
	cmds.register("agg", handlerAggregation)
	cmds.register("addfeed", middlewareLoggedIn(handlerFeedsAdd))
	cmds.register("feeds", handlerFeedsList)
	cmds.register("follow", middlewareLoggedIn(handlerFeedFollowsAdd))
	cmds.register("following", middlewareLoggedIn(handlerFeedFollowsForUser))
	cmds.register("unfollow", middlewareLoggedIn(handlerFeedFollowsDelete))

	var args []string = os.Args
	if len(args) < 2 {
		log.Fatalf("Need atleast 1 argument")
	}

	err = cmds.run(&programState, command{
		name: args[1],
		args: args[2:],
	})
	if err != nil {
		log.Fatal(err)
	}
}
