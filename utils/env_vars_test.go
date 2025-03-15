package utils

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractEmptyAttrs(t *testing.T) {
	t.Setenv("K8S_CONTAINER_NAME", "my-container")

	assert := assert.New(t)

	subject := ExtractAttrsFromEnvVar([]EnvVarEntry{})
	assert.Empty(subject)

	// invalid os env var

	entries := []EnvVarEntry{
		{"K8S_CONTAINER_NAME_FOO", "name", "process"},
	}
	subject = ExtractAttrsFromEnvVar(entries)
	assert.Empty(subject)
}

func TestExtractFromEnvVar(t *testing.T) {
	t.Setenv("K8S_CONTAINER_NAME", "my-container")
	t.Setenv("K8S_POD_NAME", "my-pod-001")
	t.Setenv("K8S_NAMESPACE", "my-ns")
	t.Setenv("K8S_SERVICE", "my-cool-svc")

	assert := assert.New(t)

	entries := []EnvVarEntry{
		{"K8S_CONTAINER_NAME", "name", "process"},
	}

	subject := entriesToAttrs(entries)
	assert.Equal("process", subject[0].Key)
	assert.Equal("Group", subject[0].Value.Kind().String())
	assert.Equal("[name=my-container]", subject[0].Value.String())
	assert.Equal("name", subject[0].Value.Resolve().Group()[0].Key)
	assert.Equal("my-container", subject[0].Value.Resolve().Group()[0].Value.String())

	// grouped env var
	entries = []EnvVarEntry{
		{"K8S_CONTAINER_NAME", "name", "process"},
		{"K8S_POD_NAME", "source", "process"},
	}

	subject = entriesToAttrs(entries)
	assert.Len(subject, 1)
	assert.Equal("Group", subject[0].Value.Kind().String())
	assert.Equal("name=my-container", subject[0].Value.Group()[0].String())
	assert.Equal("source=my-pod-001", subject[0].Value.Group()[1].String())

	// 2 groups + 1 single
	// note: group attributes are added later!
	entries = []EnvVarEntry{
		{"K8S_CONTAINER_NAME", "name", "process"},
		{"K8S_POD_NAME", "source", "process"},
		{"K8S_CONTAINER_NAME", "container_name", "kubernetes"},
		{"K8S_POD_NAME", "", "kubernetes"},
		{"K8S_SERVICE", "service", ""},
		// invalid Env name, its skipped
		{"", "container_name2", "kubernetes"},
	}

	subject = entriesToAttrs(entries)
	assert.Len(subject, 3)

	assert.Equal("my-cool-svc", subject[0].Value.String())
	assert.Equal("[name=my-container source=my-pod-001]", subject[1].Value.String())
	// no field provided, fallback to the env var name
	assert.Equal("[container_name=my-container K8S_POD_NAME=my-pod-001]", subject[2].Value.String())
}

// manual "cast" any to slog.Attr
func entriesToAttrs(entries []EnvVarEntry) []slog.Attr {
	attrs := ExtractAttrsFromEnvVar(entries)
	res := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		res[i] = a.(slog.Attr)
	}
	return res
}
