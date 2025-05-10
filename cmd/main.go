package main

import (
	"log"
	"net/http"
	"os"

	"user-pubsub/config"
	"user-pubsub/internal/handler"
	"user-pubsub/internal/queue"

	"github.com/nats-io/nats.go"
)

func main() {
	// Initialize the database connection
	if err := config.InitPostgres(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	natsURL := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("NATS connect error:", err)
	}
	defer nc.Drain()

	queue.InitConsumer(nc, "user.created")

	http.HandleFunc("/api/publish-user", handler.PublishUser(nc))
	go http.ListenAndServe(":8080", nil)

	select {} // keep service alive
}
