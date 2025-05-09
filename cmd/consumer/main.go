package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/thoraf20/fraud-detection-go/internal/stream"
	"github.com/thoraf20/fraud-detection-go/pkg/redisutil"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rdb := redisutil.NewClient("localhost:6379")
	if err := rdb.HealthCheck(ctx); err != nil {
		log.Fatalf("Redis healthcheck failed: %v", err)
	}

	consumer := &stream.Consumer{
		RDB:       rdb.Client,
		Group:     "fraud-detection",
		Consumer:  "consumer-1",
		Stream:    "transactions:raw",
		BatchSize: 10,
	}

	log.Println("Starting consumer...")
	if err := consumer.ProcessTransactions(ctx); err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
}