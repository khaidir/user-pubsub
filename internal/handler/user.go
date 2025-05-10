package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nats-io/nats.go"
)

type UserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func PublishUser(nc *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		payload, _ := json.Marshal(req)
		if err := nc.Publish("user.created", payload); err != nil {
			http.Error(w, "failed to publish", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"status":  "success",
			"message": "message published",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	}
}
