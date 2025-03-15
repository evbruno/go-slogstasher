package slogstasher

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	assert := assert.New(t)

	assert.NotNil(result["@timestamp"])
	assert.Equal("go.slogstasher", result["logger"])
	assert.Equal("1", result["version"])
	assert.Equal("INFO", result["level"])
	assert.Equal("attr", result["initial"])
	assert.Equal(uint64(42), result["u_int64"])

	source := result["source"]
	assert.NotNil(source)
}
