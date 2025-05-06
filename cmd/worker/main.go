package main

import (
	"context"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/cache"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/workers/flush"
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

	go eventbus.Consume(ctx, brokers, "auth", "analytics-svc", flush.Handle)
	go eventbus.Consume(ctx, brokers, "entries", "redis-flush", cache.HandleEntry)
	// go eventbus.Consume(ctx, brokers, "users", "something", handler)

	select {} // block forever
}
