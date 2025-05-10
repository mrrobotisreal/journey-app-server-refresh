package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	entriesworker "github.com/mrrobotisreal/journey-app-server-refresh/internal/workers/entries"
	usersworker "github.com/mrrobotisreal/journey-app-server-refresh/internal/workers/users"
	"log"
	"os"
	"strings"
)

func brokersFromEnv() []string {
	b := os.Getenv("KAFKA_BROKERS")
	if b == "" {
		b = "localhost:9092"
	}
	return strings.Split(b, ",")
}

func main() {
	_ = godotenv.Load()

	if err := db.InitDB(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	ctx := context.Background()
	brokers := brokersFromEnv()

	log.Printf("Worker is connecting to brokers: %v", brokers)

	go eventbus.Consume(ctx, brokers, "users", "users-analytics", usersworker.HandleCreateAccount)
	go eventbus.Consume(ctx, brokers, "entries", "entry-inserts", entriesworker.HandleCreateEntry)
	//go eventbus.Consume(ctx, brokers, "auth", "analytics-svc", flush.Handle)
	//go eventbus.Consume(ctx, brokers, "entries", "redis-flush", cache.HandleEntry)

	select {} // block forever
}
