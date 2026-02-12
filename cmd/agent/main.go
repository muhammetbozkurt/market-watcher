package main

import (
	"agent/internal/config"
	"log"
	"os"
	"strings"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: %w", err)
	}
	log.Println("config loaded: %v", cfg)

	log.Println("test data")
	Args := os.Args
	log.Println(strings.Join(Args[1:], "-"))
}
