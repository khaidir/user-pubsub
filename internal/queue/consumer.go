package queue

import (
	"log"

	"user-pubsub/config"
	"user-pubsub/internal/worker"

	"github.com/nats-io/nats.go"
)

func InitConsumer(nc *nats.Conn, subject string) {
	// Ensure that the database connection is properly initialized
	if config.DB == nil {
		log.Fatal("Database connection not initialized")
		return
	}

	// Create the worker handler using the GORM DB connection
	h := worker.HandlerMessage(config.DB)

	// Subscribe to the NATS subject and process messages
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		// Process each message
		if err := h.ProcessMessage(m.Data); err != nil {
			log.Printf("Error processing message: %v", err)
			return
		}
	})

	// If subscription fails, log the error
	if err != nil {
		log.Fatalf("Subscription failed: %v", err)
		return
	}
}
