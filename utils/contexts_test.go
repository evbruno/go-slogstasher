package utils

import (
	"context"
	"testing"
	"time"

	"log/slog"

	"github.com/stretchr/testify/assert"
)

type MyKey string

const (
	key1 MyKey = "key1"
	key2 MyKey = "key2"
	key3 MyKey = "key3"
)

type FakeClock struct {
	now time.Time
}

func (c FakeClock) Now() time.Time {
	return c.now
}

func TestExtractArgsFromCtx(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		ctx      context.Context
		keys     []any
		expected []slog.Attr
	}{
		{
			name: "Extract single key-value pair",
			ctx:  context.WithValue(context.Background(), key1, "value1"),
			keys: []any{key1},
			expected: []slog.Attr{
				{Key: "key1", Value: slog.StringValue("value1")},
			},
		},
		{
			name: "Extract multiple key-value pairs",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), key1, "value1")
				return context.WithValue(ctx, key2, 42)
			}(),
			keys: []any{key1, key2},
			expected: []slog.Attr{
				{Key: "key1", Value: slog.StringValue("value1")},
				{Key: "key2", Value: slog.IntValue(42)},
			},
		},
		{
			name:     "Key not found in context",
			ctx:      context.WithValue(context.Background(), key1, "value1"),
			keys:     []any{key2},
			expected: []slog.Attr{},
		},
		{
			name:     "Nil value in context",
			ctx:      context.WithValue(context.Background(), key1, nil),
			keys:     []any{"key1"},
			expected: []slog.Attr{},
		},
		{
			name: "Mixed types in context",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), key1, "value1")
				ctx = context.WithValue(ctx, key2, 42)
				return context.WithValue(ctx, key3, now)
			}(),
			keys: []any{key1, key2, key3},
			expected: []slog.Attr{
				{Key: "key1", Value: slog.StringValue("value1")},
				{Key: "key2", Value: slog.IntValue(42)},
				{Key: "key3", Value: slog.TimeValue(now)},
			},
		},
		{
			name: "ResponseTimeMsAttr empty not started",
			ctx:  context.WithValue(context.Background(), key1, "value1"),
			keys: []any{key1, StartedAtCtxKey, ResponseTimeMsKey},
			expected: []slog.Attr{
				{Key: "key1", Value: slog.StringValue("value1")},
			},
		},
		{
			name: "ResponseTimeMsAttr started no response",
			ctx: func() context.Context {
				ctx := context.WithValue(context.Background(), key1, "value1")
				return context.WithValue(ctx, StartedAtCtxKey, now)
			}(),
			keys: []any{key1, StartedAtCtxKey, ResponseTimeMsKey},
			expected: []slog.Attr{
				{Key: "key1", Value: slog.StringValue("value1")},
				{Key: "startedAt", Value: slog.TimeValue(now)},
			},
		},
		{
			name: "ResponseTimeMsAttr started with response",
			ctx: func() context.Context {
				myClock0 := FakeClock{now}
				myClock1 := FakeClock{now: now.Add(42 * time.Millisecond)}
				ctx := context.WithValue(context.Background(), key1, "value1")
				ctx = WithStartedNowCtx(ctx, myClock0)
				return WithResponseTimeMsCtx(ctx, myClock1)
			}(),
			keys: []any{key1, StartedAtCtxKey, ResponseTimeMsKey},
			expected: []slog.Attr{
				{Key: "key1", Value: slog.StringValue("value1")},
				{Key: "startedAt", Value: slog.TimeValue(now)},
				{Key: "responseTimeMs", Value: slog.Int64Value(42)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractArgsFromCtx(tt.ctx, tt.keys...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
