package main

import (
	"log"

	"github.com/adal4ik/people-enrichment-service/internal/logger"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
}
