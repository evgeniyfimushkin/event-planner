package middlewarelogger

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareLogger(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	middlewareLogger := New(logger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrappedHandler := middlewareLogger(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "GoTest")
	req.RemoteAddr = "127.0.0.1:8080"

	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, `"component":"middleware/logger"`)
	assert.Contains(t, logOutput, `"method":"GET"`)
	assert.Contains(t, logOutput, `"path":"/test"`)
	assert.Contains(t, logOutput, `"remote_addr":"127.0.0.1:8080"`)
	assert.Contains(t, logOutput, `"user_agent":"GoTest"`)
	assert.Contains(t, logOutput, `"request_id":"test-request-id"`)
	assert.Contains(t, logOutput, `"status":200`)
}

