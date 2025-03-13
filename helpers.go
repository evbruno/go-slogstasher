package slogstasher

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
)

func (h *Logtsash) formatMessage(_ context.Context, record *slog.Record) map[string]any {
	log := map[string]any{
		"@timestamp": record.Time.UTC(),
		"logger":     "go.slogstasher", // FIXME: hardcoded
		"version":    "1",              // FIXME: hardcoded
		"level":      record.Level.String(),
		"message":    record.Message,
	}

	allAttrs := h.attrs

	if h.opts.AddSource {
		allAttrs = append(allAttrs, addSource("source", record))
	}

	if h.opts.ReplaceAttr != nil {
		for i, attr := range allAttrs {
			allAttrs[i] = h.opts.ReplaceAttr(h.groups, attr)
		}
	}

	for _, attr := range allAttrs {
		k, v := formatAttrs(&attr)
		log[k] = v
	}

	return log
}

func formatAttrs(attr *slog.Attr) (string, any) {
	v := attr.Value
	kind := v.Kind()

	switch kind {
	case slog.KindString:
		return attr.Key, v.String()
	case slog.KindAny:
		return attr.Key, v.Any()
	case slog.KindTime:
		return attr.Key, v.Time()
	case slog.KindUint64:
		return attr.Key, v.Uint64()
	case slog.KindInt64:
		return attr.Key, v.Int64()
	case slog.KindFloat64:
		return attr.Key, v.Float64()
	case slog.KindBool:
		return attr.Key, v.Bool()
	case slog.KindDuration:
		return attr.Key, v.Duration()
	case slog.KindLogValuer:
		return attr.Key, v.Any()
	case slog.KindGroup:
		children := map[string]any{}
		for _, a := range v.Group() {
			ck, cv := formatAttrs(&a)
			children[ck] = cv
		}
		return attr.Key, children

	default:
		// TODO: handle error
		panic(fmt.Sprintf("unsupported kind: %s", kind))
	}

}

// thanks Gepetto  🤖🤖
func addSource(sourceKey string, r *slog.Record) slog.Attr {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	var args []any
	if f.Function != "" {
		args = append(args, slog.String("function", f.Function))
	}
	if f.File != "" {
		args = append(args, slog.String("file", f.File))
	}
	if f.Line != 0 {
		args = append(args, slog.Int("line", f.Line))
	}

	return slog.Group(sourceKey, args...)
}