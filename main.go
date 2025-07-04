package main

import (
	"log"
	"os"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	log.Printf("Blog Aggregator")

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config Read Error: %v", err)
	}
	programState := state{
		cfg: &cfg,
	}
	cmds := commands{
		nameToFunction: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	var args []string = os.Args
	if len(args) < 2 {
		log.Fatalf("Need atleast 1 argument")
	}
	com := command{
		name: args[1],
		args: args[2:],
	}
	err = cmds.run(&programState, com)
	if err != nil {
		log.Fatal(err)
	}
}
