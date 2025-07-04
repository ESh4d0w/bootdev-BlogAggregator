package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/config"
	"github.com/esh4d0w/bootdev-BlogAggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	log.Printf("Blog Aggregator")

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config Read Error: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	dbQueries := database.New(db)

	programState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		nameToFunction: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUserList)

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
