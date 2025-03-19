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

func TestExtractSingleKeyValuePair(t *testing.T) {
	ctx := context.WithValue(context.Background(), key1, "value1")
	keys := []any{key1}

	expected := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
	}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestExtractMultipleKeyValuePairs(t *testing.T) {
	ctx := context.WithValue(context.Background(), key1, "value1")
	ctx = context.WithValue(ctx, key2, 42)
	keys := []any{key1, key2}

	expected := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
		{Key: "key2", Value: slog.IntValue(42)},
	}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestKeyNotFoundInContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), key1, "value1")
	keys := []any{key2}

	expected := []slog.Attr{}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestNilValueInContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), key1, nil)
	keys := []any{"key1"}

	expected := []slog.Attr{}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestMixedTypesInContext(t *testing.T) {
	now := time.Now()
	ctx := context.WithValue(context.Background(), key1, "value1")
	ctx = context.WithValue(ctx, key2, 42)
	ctx = context.WithValue(ctx, key3, now)
	keys := []any{key1, key2, key3}

	expected := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
		{Key: "key2", Value: slog.IntValue(42)},
		{Key: "key3", Value: slog.TimeValue(now)},
	}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestResponseTimeMsAttrEmptyNotStarted(t *testing.T) {
	ctx := context.WithValue(context.Background(), key1, "value1")
	keys := []any{key1, StartedAtCtxKey, ResponseTimeMsKey}

	expected := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
	}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestResponseTimeMsAttrStartedNoResponse(t *testing.T) {
	now := time.Now()

	ctx := context.WithValue(context.Background(), key1, "value1")
	ctx = context.WithValue(ctx, StartedAtCtxKey, now)

	keys := []any{key1, StartedAtCtxKey, ResponseTimeMsKey}

	expected := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
		{Key: "startedAt", Value: slog.TimeValue(now)},
	}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}

func TestResponseTimeMsAttrStartedWithResponse(t *testing.T) {
	now := time.Now()

	myClock0 := FakeClock{now}

	myClock1 := FakeClock{
		now: now.Add(42 * time.Millisecond),
	}

	ctx := context.WithValue(context.Background(), key1, "value1")
	ctx = WithStartedNowCtx(ctx, myClock0)
	ctx = WithResponseTimeMsCtx(ctx, myClock1)

	keys := []any{key1, StartedAtCtxKey, ResponseTimeMsKey}

	expected := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
		{Key: "startedAt", Value: slog.TimeValue(now)},
		{Key: "responseTimeMs", Value: slog.Int64Value(42)},
	}

	result := ExtractArgsFromCtx(ctx, keys...)
	assert.Equal(t, expected, result)
}
