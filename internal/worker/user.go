package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"user-pubsub/model"

	"github.com/jinzhu/gorm" // Import GORM
)

type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new Handler instance with the provided GORM DB connection
func HandlerMessage(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// ProcessMessage handles the processing of incoming user messages
func (h *Handler) ProcessMessage(data []byte) error {
	if h.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		log.Printf("failed to unmarshal message: %v", err)
		return fmt.Errorf("invalid message format: %w", err)
	}

	// Using WaitGroup to handle goroutine synchronization for concurrent processing
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done() // Decrement the counter when the goroutine is finished

		// Insert the user into the database using GORM
		if err := h.db.Create(&user).Error; err != nil {
			log.Printf("database insert failed: %v", err)
			return
		}
	}()

	// Wait for all goroutines to finish before returning
	wg.Wait()

	return nil
}
