[![Build status](https://img.shields.io/github/actions/workflow/status/evbruno/go-slogstasher/build-and-test.yml?style=for-the-badge&branch=main)](https://github.com/evbruno/go-slogstasher/actions/actions?workflow=build-and-test)
[![Release](https://img.shields.io/github/release/evbruno/go-slogstasher.svg?style=for-the-badge)](https://github.com/evbruno/go-slogstasher/releases/latest)

# Go slog.Handler for Logstash

**work in progress**

## Usage

Start a TCP server in one shell:

```bash
socat -u TCP-LISTEN:4560,reuseaddr,fork STDOUT
```

## Go

Simple usage:

```go

import (
	"log/slog"

	s "github.com/evbruno/go-slogstasher"
)

stash := &s.LogstashOpts{
	Host:  "127.0.0.1",
	Port:  4560,
	Type:  s.Tcp4,
}

slog.SetDefault(slog.New(s.NewLogstashHandler(stash)))

slog.Info("Hello world of Go!")
```

Outputs:

```
{"@timestamp":"2025-03-13T13:34:52.514582Z","level":"INFO","logger":"go.slogstasher","message":"Hello world of Go!","version":"1"}
```

## Utils

Load a set of fields from the env var, and return the appropiate `slog.Attr`.

**Note**: type is `[]any` to make our life easier when dealing with `slog.With`.

```go

// import su "github.com/evbruno/go-slogstasher/utils"

envVarAttrs :=  []su.EnvVarEntry{
	{Key: "K8S_CONTAINER_NAME", Attr: "name", Group: "process"},
	{Key: "K8S_POD_NAME", Attr: "source", Group: "process"},
	{Key: "K8S_SERVICE", Attr: "service"},
}

newAttrs []any = su.ExtractAttrsFromEnvVar(envVarAttrs)

// use it!

log.With(attrs...).Info("Hi there env vars !")

```

Outputs:

```
2025/03/14 20:56:56 INFO Hi there env vars ! service=my-service kubernetes.container=my-container kubernetes.pod=my-pod-0001
```


