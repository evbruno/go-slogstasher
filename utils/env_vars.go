package utils

import (
	"log/slog"
	"os"
)

type EnvVarEntry struct {
	Key   string
	Attr  string
	Group string
}

func ExtractAttrsFromEnvVar(fields []EnvVarEntry) []any {
	var attrs []any
	var groups map[string][]any = make(map[string][]any)

	for _, f := range fields {
		if f.Key == "" {
			continue
		}

		if value := os.Getenv(f.Key); value != "" {
			fieldName := f.Attr
			if fieldName == "" {
				fieldName = f.Key
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
