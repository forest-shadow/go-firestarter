package httpserver

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/forest-shadow/go-firestarter/pkg/logger"
)

func TestNewBuildsConfiguredServer(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Port:              "9090",
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       4 * time.Second,
		ShutdownTimeout:   5 * time.Second,
	}

	server, err := New(http.NotFoundHandler(), cfg, newTestLogger())
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	if server.server.Addr != ":9090" {
		t.Fatalf("expected server address %q, got %q", ":9090", server.server.Addr)
	}

	if server.server.ReadTimeout != time.Second {
		t.Fatalf("expected read timeout %s, got %s", time.Second, server.server.ReadTimeout)
	}

	if server.server.ReadHeaderTimeout != 2*time.Second {
		t.Fatalf("expected read header timeout %s, got %s", 2*time.Second, server.server.ReadHeaderTimeout)
	}

	if server.server.WriteTimeout != 3*time.Second {
		t.Fatalf("expected write timeout %s, got %s", 3*time.Second, server.server.WriteTimeout)
	}

	if server.server.IdleTimeout != 4*time.Second {
		t.Fatalf("expected idle timeout %s, got %s", 4*time.Second, server.server.IdleTimeout)
	}
}

func TestNewRejectsInvalidDependenciesAndConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler http.Handler
		config  Config
		log     logger.Logger
		want    string
	}{
		{
			name:    "invalid config",
			handler: http.NotFoundHandler(),
			config:  Config{ReadTimeout: -time.Second},
			log:     newTestLogger(),
			want:    "http_server.read_timeout",
		},
		{
			name:   "missing handler",
			config: Config{},
			log:    newTestLogger(),
			want:   "handler is required",
		},
		{
			name:    "missing logger",
			handler: http.NotFoundHandler(),
			config:  Config{},
			want:    "logger is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tt.handler, tt.config, tt.log)
			if err == nil {
				t.Fatal("expected error")
			}

			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %q", tt.want, err)
			}
		})
	}
}

func TestRunReturnsListenError(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	defer listener.Close()

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	server, err := New(http.NotFoundHandler(), Config{Port: port}, newTestLogger())
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	if err := server.Run(); err == nil {
		t.Fatal("expected listen error")
	}
}

func TestRunAndShutdown(t *testing.T) {
	t.Parallel()

	log := newTestLogger()
	server, err := New(http.NotFoundHandler(), Config{Port: "0"}, log)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	runErr := make(chan error, 1)
	go func() {
		runErr <- server.Run()
	}()

	select {
	case msg := <-log.info:
		if msg != "http server: started" {
			t.Fatalf("expected started log, got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not start")
	}

	if err := server.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown server: %v", err)
	}

	select {
	case err := <-runErr:
		if err != nil {
			t.Fatalf("run server: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not stop")
	}
}

type testLogger struct {
	info chan string
}

func newTestLogger() *testLogger {
	return &testLogger{info: make(chan string, 2)}
}

func (*testLogger) Debug(string, ...logger.Field) {}

func (l *testLogger) Info(msg string, _ ...logger.Field) {
	l.info <- msg
}

func (*testLogger) Warn(string, ...logger.Field) {}

func (*testLogger) Error(string, ...logger.Field) {}

func (l *testLogger) With(...logger.Field) logger.Logger {
	return l
}

func (*testLogger) Sync() error {
	return nil
}
