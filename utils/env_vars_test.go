package utils

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAttrsFromEnvVar(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		entries    []EnvVarEntry
		expected   []string
		expectedLen int
	}{
		{
			name:       "Empty attributes",
			envVars:    map[string]string{"K8S_CONTAINER_NAME": "my-container"},
			entries:    []EnvVarEntry{},
			expected:   []string{},
			expectedLen: 0,
		},
		{
			name:       "Invalid OS env var",
			envVars:    map[string]string{"K8S_CONTAINER_NAME": "my-container"},
			entries:    []EnvVarEntry{{"K8S_CONTAINER_NAME_FOO", "name", "process"}},
			expected:   []string{},
			expectedLen: 0,
		},
		{
			name: "Single group attribute",
			envVars: map[string]string{
				"K8S_CONTAINER_NAME": "my-container",
			},
			entries: []EnvVarEntry{
				{"K8S_CONTAINER_NAME", "name", "process"},
			},
			expected:   []string{"[name=my-container]"},
			expectedLen: 1,
		},
		{
			name: "Grouped env vars",
			envVars: map[string]string{
				"K8S_CONTAINER_NAME": "my-container",
				"K8S_POD_NAME":       "my-pod-001",
			},
			entries: []EnvVarEntry{
				{"K8S_CONTAINER_NAME", "name", "process"},
				{"K8S_POD_NAME", "source", "process"},
			},
			expected:   []string{"[name=my-container source=my-pod-001]"},
			expectedLen: 1,
		},
		{
			name: "Multiple groups and single attributes",
			envVars: map[string]string{
				"K8S_CONTAINER_NAME": "my-container",
				"K8S_POD_NAME":       "my-pod-001",
				"K8S_SERVICE":        "my-cool-svc",
			},
			entries: []EnvVarEntry{
				{"K8S_CONTAINER_NAME", "name", "process"},
				{"K8S_POD_NAME", "source", "process"},
				{"K8S_CONTAINER_NAME", "container_name", "kubernetes"},
				{"K8S_POD_NAME", "", "kubernetes"},
				{"K8S_SERVICE", "service", ""},
				{"", "container_name2", "kubernetes"}, // Invalid entry
			},
			expected: []string{
				"my-cool-svc",
				"[name=my-container source=my-pod-001]",
				"[container_name=my-container K8S_POD_NAME=my-pod-001]",
			},
			expectedLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			subject := ExtractAttrsFromEnvVar(tt.entries)
			assert.Len(t, subject, tt.expectedLen)

			sort.Slice(subject, func(i, j int) bool {
				return subject[i].Key > subject[j].Key
			})

			for i, expected := range tt.expected {
				assert.Equal(t, expected, subject[i].Value.String())
			}
		})
	}
}
