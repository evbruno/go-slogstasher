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

Load a set of fields from the env var, and return the appropiate `slog.Attr` (or `any`).

```go

// import su "github.com/evbruno/go-slogstasher/utils"

envVarAttrs :=  []su.EnvVarEntry{
	{Env: "K8S_CONTAINER_NAME", Attr: "name", Group: "process"},
	{Env: "K8S_POD_NAME", Attr: "source", Group: "process"},
	{Env: "K8S_SERVICE", Attr: "service"},
}

var args []any = su.ExtractArgsFromEnvVar(opts)
slog.With(args...).Info("Hi there env vars + args")

```

Outputs:

```
2025/03/14 21:18:04 INFO Hi there env vars + args service=my-service kubernetes.container=my-container kubernetes.pod=my-pod-0001
```

*... or use the alternative method:*

```go

var attrs []slog.Attr = su.ExtractAttrsFromEnvVar(opts)
logger := slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs(attrs))
logger.Info("Hi there env vars with attrs !")

```

Outputs:

```
time=2025-03-14T21:18:04.831-03:00 level=INFO msg="Hi there env vars with attrs !" service=my-service kubernetes.container=my-container kubernetes.pod=my-pod-0001
```