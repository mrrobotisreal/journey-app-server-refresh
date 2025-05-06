package main

import (
	"context"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
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
	ctx := context.Background()
	brokers := brokersFromEnv()

	log.Printf("Worker is connecting to brokers: %v", brokers)

	go eventbus.Consume(ctx, brokers, "users", "users-analytics", usersworker.HandleCreateAccount)
	//go eventbus.Consume(ctx, brokers, "auth", "analytics-svc", flush.Handle)
	//go eventbus.Consume(ctx, brokers, "entries", "redis-flush", cache.HandleEntry)

	select {} // block forever
}
