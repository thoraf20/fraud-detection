package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

type Transaction struct {
	ID        string  `json:"transaction_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Merchant  string  `json:"merchant"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	// Redis client
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	merchants := []string{"Amazon", "Jumia", "Tesla", "FraudulentFoods", "Temu"}

	for {
		// Simulate a new transaction
		txn := Transaction{
			ID:        "txn_" + randString(8),
			UserID:    "user_" + randString(4),
			Amount:    rand.Float64() * 1000,
			Merchant:  merchants[rand.Intn(len(merchants))],
			Timestamp: time.Now().Unix(),
		}

		// Serialize to JSON
		values, _ := json.Marshal(txn)

		// Push to Redis Stream
		_, err := rdb.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "transactions:raw",
			Values: map[string]interface{}{"data": string(values)},
		}).Result()

		if err != nil {
			log.Printf("Failed to publish transaction: %v", err)
		} else {
			log.Printf("Produced transaction: %s", txn.ID)
		}

		time.Sleep(500 * time.Millisecond) // Simulate realistic throughput
	}
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}