package main

import (
	"context"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/workers/flush"
)

func main() {
	ctx := context.Background()
	brokers := []string{"kafka:9092"}

	//go eventbus.Consume(ctx, brokers, "auth", "analytics-svc", analytics.HandleLogin)
	//go eventbus.Consume(ctx, brokers, "entries", "analytics-svc", analytics.HandleEntry)
	go eventbus.Consume(ctx, brokers, "entries", "redis-flush", flush.Handle)

	select {} // block forever
}
