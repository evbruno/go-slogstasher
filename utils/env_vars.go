package utils

import (
	"log/slog"
	"os"
)

type EnvVarEntry struct {
	Env   string
	Attr  string
	Group string
}

func ExtractArgsFromEnvVar(entries []EnvVarEntry) []any {
	attrs := ExtractAttrsFromEnvVar(entries)
	res := make([]any, len(attrs))

	for i, a := range attrs {
		res[i] = a
	}

	return res
}

func ExtractAttrsFromEnvVar(entries []EnvVarEntry) []slog.Attr {
	var attrs []slog.Attr
	var groups map[string][]any = make(map[string][]any)

	for _, f := range entries {
		if f.Env == "" {
			continue
		}

		if value := os.Getenv(f.Env); value != "" {
			fieldName := f.Attr
			if fieldName == "" {
				fieldName = f.Env
			}

			if f.Group != "" {
				groups[f.Group] = append(groups[f.Group], slog.String(fieldName, value))
			} else {
				attrs = append(attrs, slog.String(fieldName, value))
			}
		}
	}

	for k, v := range groups {
		attrs = append(attrs, slog.Group(k, v...))
	}

	return attrs
}
