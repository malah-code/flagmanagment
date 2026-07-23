package logging

import (
	"bytes"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name      string
		logFormat string
		env       string
		wantJSON  bool
	}{
		{"text development", "text", "development", false},
		{"json production", "json", "production", true},
		{"auto development", "auto", "development", false},
		{"auto production", "auto", "production", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(tt.logFormat, tt.env, &buf)

			logger.Info().Msg("test message")

			output := buf.String()
			isJSON := output[0] == '{'

			if isJSON != tt.wantJSON {
				t.Errorf("expected JSON output: %v, got: %v (output: %s)", tt.wantJSON, isJSON, output)
			}
		})
	}
}
