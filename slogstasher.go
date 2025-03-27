package slogstasher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
)

// public types/consts

type ConnectionType string

const (
	Udp  ConnectionType = "udp"
	Tcp  ConnectionType = "tcp"
	Tcp4 ConnectionType = "tcp4"
)

type LogstashOpts struct {
	Host  string
	Port  int
	Type  ConnectionType
	Level slog.Leveler

	Conn net.Conn

	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

type Logtsash struct {
	opts   *LogstashOpts
	attrs  []slog.Attr
	groups []string
	conn   net.Conn
}

// API

func NewLogstashHandler(opts *LogstashOpts) slog.Handler {
	if opts.Conn == nil {
		addr := net.JoinHostPort(opts.Host, fmt.Sprintf("%d", opts.Port))
		// FIMXE: DEBUG?
		fmt.Println("[DEBUG] Connecting to logstash at", addr, string(opts.Type))
		conn, err := net.Dial(string(opts.Type), addr)

		if err != nil {
			// DEBUG?
			fmt.Println("[DEBUG] Error connecting to logstash, fallback to stdout, err:", err)
			return nil
		}
		opts.Conn = conn
	}

	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}

	return &Logtsash{
		conn:   opts.Conn,
		opts:   opts,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

// Handler interface

func (h *Logtsash) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Logtsash) Handle(ctx context.Context, record slog.Record) error {
	payload := h.formatMessage(ctx, &record)
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	go func() {
		_, err = h.conn.Write(append(bytes, byte('\n')))
		//FIXME: DEBUG? retry?
		if err != nil {
			fmt.Println("Error writing to logstash:", err, string(bytes))
		} else {
			fmt.Println("[DEBUG]", string(bytes))
		}
	}()

	return err
}

func (h *Logtsash) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Logtsash{
		conn:   h.conn,
		opts:   h.opts,
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *Logtsash) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &Logtsash{
		conn:   h.conn,
		opts:   h.opts,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
