package utils

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type CommonCtxKey string

const (
	StartedAtCtxKey   CommonCtxKey = "startedAt"
	ResponseTimeMsKey CommonCtxKey = "responseTimeMs"
)

// ContextualizedHandler is a custom implementation of slog.Handler that allows
// for adding contextual keys to log entries. It embeds the slog.Handler interface
// and includes a slice of keys (`Keys`) that can be used to provide additional
// context to log messages.
type ContextualizedHandler struct {
	slog.Handler
	Keys []any
}

func (h *ContextualizedHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, attr := range ExtractArgsFromCtx(ctx, h.Keys...) {
		r.AddAttrs(attr)
	}

	return h.Handler.Handle(ctx, r)
}

// Clock is an interface that provides an abstraction for retrieving the current time.
// It defines a single method, Now, which returns the current time as a time.Time value.
//
// This interface is particularly useful for testing purposes, as it allows you to
// mock or substitute the implementation of time retrieval. For example, you can
// use a custom implementation of Clock to simulate specific times or control
// the flow of time in your tests.
type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func WithStartedNowCtx(ctx context.Context, c Clock) context.Context {
	var now time.Time = c.Now()
	return context.WithValue(ctx, StartedAtCtxKey, now)
}

func ResponseTimeMsAttr(ctx context.Context, c Clock) *slog.Attr {
	if startedAt, ok := ctx.Value(StartedAtCtxKey).(time.Time); ok {
		value := c.Now().Sub(startedAt).Milliseconds()
		attr := slog.Int64(string(ResponseTimeMsKey), value)
		return &attr
	}

	return nil
}

func WithResponseTimeMsCtx(ctx context.Context, c Clock) context.Context {
	if startedAt, ok := ctx.Value(StartedAtCtxKey).(time.Time); ok {
		//value := time.Since(startedAt).Milliseconds()
		value := c.Now().Sub(startedAt).Milliseconds()
		return context.WithValue(ctx, ResponseTimeMsKey, value)
	}

	return ctx
}

// ExtractArgsFromCtx extracts a slice of slog.Attr from the provided context
// based on the given keys. For each key, it retrieves the corresponding value
// from the context and converts it into a slog.Attr using the anyToAttr function.
// If the value is nil or cannot be converted, it is skipped.
//
// Parameters:
//   - ctx: The context from which to extract values.
//   - keys: A variadic list of keys to look up in the context.
//
// Returns:
//
//	A slice of slog.Attr containing the extracted attributes.
func ExtractArgsFromCtx(ctx context.Context, keys ...any) []slog.Attr {
	attrs := []slog.Attr{}
	for _, key := range keys {
		if value := ctx.Value(key); value != nil {
			//if attr := anyToAttr(fmt.Sprintf("%v", key), value); attr != nil {
			attr := anyToAttr(fmt.Sprintf("%v", key), value)
			attrs = append(attrs, attr)
			//}
		}
	}
	return attrs
}

func anyToAttr(key string, variable any) slog.Attr {
	switch v := variable.(type) {
	case int:
		return slog.Attr{Key: key, Value: slog.IntValue(v)}
	case int64:
		return slog.Attr{Key: key, Value: slog.Int64Value(v)}
	case float32:
		return slog.Attr{Key: key, Value: slog.Float64Value(float64(v))}
	case float64:
		return slog.Attr{Key: key, Value: slog.Float64Value(v)}
	case string:
		return slog.Attr{Key: key, Value: slog.StringValue(v)}
	case bool:
		return slog.Attr{Key: key, Value: slog.BoolValue(v)}
	case time.Time:
		return slog.Attr{Key: key, Value: slog.TimeValue(v)}
	default:
		// TODO: add more types or `Any` is good enough?
		return slog.Attr{Key: key, Value: slog.AnyValue(v)}
	}
}
