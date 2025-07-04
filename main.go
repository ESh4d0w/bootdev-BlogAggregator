package main

import (
	"fmt"
	"log"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/config"
)

func main() {
	fmt.Println("Blog Aggregator")

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Config Read Error: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("ESh4d0w")
	if err != nil {
		log.Fatalf("Error setting User: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Config ReRead Error: %v", err)
	}

	fmt.Printf("Read config again: %+v\n", cfg)
}
