package stream

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Consumer struct {
	RDB       *redis.Client
	Group     string
	Consumer  string
	Stream    string
	BatchSize int64
}

func (c *Consumer) ProcessTransactions(ctx context.Context) error {
	_, err := c.RDB.XGroupCreateMkStream(ctx, c.Stream, c.Group, "$").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}

	for {
		entries, err := c.RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.Group,
			Consumer: c.Consumer,
			Streams:  []string{c.Stream, ">"},
			Count:    c.BatchSize,
			Block:    200 * time.Millisecond,
		}).Result()

		if err != nil && err != redis.Nil {
			return err
		}

		// Process each message
		for _, entry := range entries[0].Messages {
			var txn Transaction
			if err := json.Unmarshal([]byte(entry.Values["data"].(string)), &txn); err != nil {
				log.Printf("Failed to unmarshal transaction: %v", err)
				continue
			}

			log.Printf("Consumed transaction: %s (Amount: $%.2f)", txn.ID, txn.Amount)

			// TODO: Add fraud detection logic here
			// ...

			// Acknowledge processing
			c.RDB.XAck(ctx, c.Stream, c.Group, entry.ID)
		}
	}
}