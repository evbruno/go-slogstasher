package slogstasher

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestFormatMessage(t *testing.T) {
	ctx := context.Background()
	record := &slog.Record{
		Time:    time.Now(),
		Level:   slog.LevelInfo,
		Message: "test message",
	}

	opts := &LogstashOpts{
		AddSource: true,
	}

	h := &Logtsash{
		attrs: []slog.Attr{
			slog.String("initial", "attr"),
			slog.Uint64("u_int64", 42),
		},
		opts: opts,
	}

	result := h.formatMessage(ctx, record)

	if result["@timestamp"] == nil {
		t.Errorf("Expected @timestamp to be set")
	}
	if result["logger"] != "go.slogstasher" {
		t.Errorf("Expected logger to be 'go.slogstasher', got %v", result["logger"])
	}
	if result["version"] != "1" {
		t.Errorf("Expected version to be '1', got %v", result["version"])
	}
	if result["level"] != slog.LevelInfo.String() {
		t.Errorf("Expected level to be %v, got %v", slog.LevelInfo.String(), result["level"])
	}
	if result["message"] != "test message" {
		t.Errorf("Expected message to be 'test message', got %v", result["message"])
	}
	if result["initial"] != "attr" {
		t.Errorf("Expected initial to be 'attr', got %v", result["initial"])
	}
	if result["u_int64"] != uint64(42) {
		t.Errorf("Expected u_int64 to be 64, got %v", result["u_int64"])
	}
	if _, ok := result["source"]; !ok {
		t.Errorf("Expected source to be set")
	}
}
