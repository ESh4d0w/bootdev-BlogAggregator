package main

import (
	"fmt"
	"log"

	"github.com/esh4d0w/bootdev-BlogAggregator/internal/config"
)

func main() {
	fmt.Println("Blog Aggregator")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}
	err = cfg.SetUser("ESh4d0w")
	if err != nil {
		log.Fatalf("Error setting User: %v", err)
	}
	fmt.Printf("%+v", cfg)

}
