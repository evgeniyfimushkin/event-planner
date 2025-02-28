package logger

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		env           string
		expectedLevel slog.Level
	}{
		{"local", slog.LevelDebug},
		{"dev", slog.LevelDebug},
		{"prod", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			log := SetupLogger(tt.env)
			assert.NotNil(t, log, "Logger should not be nil")

			log.Info("test log message")

			w.Close()
			os.Stdout = oldStdout

			out := make([]byte, 1024)
			n, _ := r.Read(out)
			logOutput := string(out[:n])

			assert.Contains(t, logOutput, "test log message", "Logger should output the test message")
		})
	}
}

