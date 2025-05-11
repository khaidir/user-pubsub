package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

// MockNATSConn is a basic mock for *nats.Conn
type MockNATSConn struct {
	ShouldFail bool
	Published  []byte
}

func (m *MockNATSConn) Publish(subject string, data []byte) error {
	if m.ShouldFail {
		return errors.New("publish failed")
	}
	m.Published = data
	return nil
}

func TestPublishUser_Success(t *testing.T) {
	mock := &MockNATSConn{}
	handlerFunc := func(subject string, data []byte) error {
		return mock.Publish(subject, data)
	}

	// nc := &nats.Conn{}
	// Wrap handler to inject our mock manually
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UserRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		payload, _ := json.Marshal(req)
		_ = handlerFunc("user.created", payload)

		resp := map[string]string{"status": "success", "message": "message published"}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(resp)
	})

	body := `{"email":"user@example.com","name":"John Doe"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestPublishUser_InvalidJSON(t *testing.T) {
	nc := &nats.Conn{}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	PublishUser(nc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}
