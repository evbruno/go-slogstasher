![Go Tests](https://github.com/evbruno/go-slogstasher/actions/workflows/build-and-test.yml/badge.svg)

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

	s "github.com/evbruno/go-slogstasher/v1"
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



