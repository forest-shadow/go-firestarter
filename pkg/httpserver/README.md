# `pkg/httpserver`

## Purpose

`pkg/httpserver` wraps the standard `net/http` server setup used by this
starter.

It currently owns:

- HTTP listen address construction
- server read, write, and idle timeouts
- graceful shutdown timeout
- startup and shutdown logging

## Configuration

`Config` uses `time.Duration` for timeout fields:

```go
type Config struct {
 Port              string        `mapstructure:"port"`
 ReadTimeout       time.Duration `mapstructure:"read_timeout"`
 ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
 WriteTimeout      time.Duration `mapstructure:"write_timeout"`
 IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
 ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
}
```

YAML config should use Go duration strings:

```yaml
http_server:
  port: "8080"
  read_timeout: 20s
  read_header_timeout: 10s
  write_timeout: 20s
  idle_timeout: 60s
  shutdown_timeout: 25s
```

The same values can be provided through environment variables or secret-backed
environment injection:

```text
APP_HTTP_SERVER_PORT=8080
APP_HTTP_SERVER_READ_TIMEOUT=20s
APP_HTTP_SERVER_READ_HEADER_TIMEOUT=10s
APP_HTTP_SERVER_WRITE_TIMEOUT=20s
APP_HTTP_SERVER_IDLE_TIMEOUT=60s
APP_HTTP_SERVER_SHUTDOWN_TIMEOUT=25s
```

This works through the shared `pkg/config.DecodeHook()`, which enables
`mapstructure` duration parsing from string values.

## Defaults and validation

`Config.WithDefaults()` applies practical local defaults:

- port: `8080`
- read timeout: `20s`
- read header timeout: `10s`
- write timeout: `20s`
- idle timeout: `60s`
- shutdown timeout: `25s`

`Config.Validate()` ensures that the port is present and all timeouts are
positive before the server is used.

## Lifecycle API

`New()` applies defaults, validates the configuration and dependencies, and
constructs the underlying `net/http.Server`. It does not open a listener or
start a goroutine.

`Run()` opens the configured TCP listener and serves requests until shutdown or
an error occurs. It blocks the calling goroutine and returns listener and server
errors to the caller. This lets the application decide how to coordinate the
server with signal handling and other runtime components.

`Shutdown(ctx)` performs graceful shutdown and returns any shutdown error. The
configured `shutdown_timeout` is applied as an upper bound in addition to the
caller's context deadline.

```go
server, err := httpserver.New(handler, cfg, log)
if err != nil {
 return err
}

runErr := make(chan error, 1)
go func() {
 runErr <- server.Run()
}()

// Application-level lifecycle code waits for a signal or runErr.

if err := server.Shutdown(ctx); err != nil {
 return err
}
```

The package intentionally does not own OS signal handling or application-level
goroutine coordination.

## `HTTP_READ_TIMEOUT`

### Purpose

> `HTTP_READ_TIMEOUT` defines the maximum amount of time the server allows for reading an incoming HTTP request.

It helps protect the server from hanging clients, slow connections, and slow request-body attacks.

This timeout is useful when a client opens a connection but does not complete the request in a reasonable time.

Examples:

- the client starts sending a request but hangs;
- the client has an unstable or very slow network connection;
- the client sends the request body extremely slowly;
- an attacker intentionally keeps many connections open to exhaust server resources.

### Why it matters

Each open connection consumes server resources, such as:

- a socket;
- memory;
- a goroutine or worker;
- connection limits;
- internal queues or request-handling capacity.

If many slow or hanging clients keep connections open, normal users may experience increased latency, errors, or service unavailability.

### Behavior

If the client does not finish sending the request within the configured timeout, the server closes the connection.

In simple terms: the server is willing to wait for a reasonable amount of time, but not forever.

### Recommended values

For most HTTP APIs, a good default is:

- `5s–10s` for regular APIs with small requests;
- `15s–30s` for slower clients, mobile networks, or larger request bodies;
- `60s+` only for special cases, such as large uploads.

Longer timeouts should be used carefully because they allow slow clients to hold server resources for a longer time.

## `HTTP_READ_HEADER_TIMEOUT`

### Purpose

> `HTTP_READ_HEADER_TIMEOUT` defines the maximum amount of time the server allows for reading HTTP request headers.

It protects the server from clients that open a connection but send the request line or headers too slowly.

### What it protects against

This timeout is mainly useful against slow header delivery and Slowloris-style behavior.

> **Slowloris** is a denial-of-service attack where an attacker opens many HTTP connections and keeps them alive by sending incomplete request headers very slowly, forcing the server to waste resources waiting for the requests to finish.

Examples:
- the client opens a TCP connection but delays sending HTTP headers;
- the client sends headers byte by byte very slowly;
- the client starts a request but never completes the header section;
- an attacker keeps many connections open by slowly streaming incomplete headers.

### Why it matters

Before the server receives complete headers, it usually cannot route the request to the correct handler or make request-specific decisions.

Each incomplete request may still consume server resources, such as:
- a socket;
- memory;
- a goroutine or worker;
- connection slots;
- internal server capacity.

If many clients keep connections stuck at the header-reading stage, the server may become unavailable for normal users.

### Behavior

If the client does not finish sending request headers within the configured timeout, the server closes the connection.

In simple terms: the server gives the client a short window to send valid HTTP headers, but it does not wait forever.

### Relationship to `HTTP_READ_TIMEOUT`

- `HTTP_READ_HEADER_TIMEOUT` is narrower than `HTTP_READ_TIMEOUT`.
- `HTTP_READ_HEADER_TIMEOUT` applies only to reading the request line and headers.
- `HTTP_READ_TIMEOUT` applies to reading the entire request, including headers and the request body.

This distinction is important:
- use `HTTP_READ_HEADER_TIMEOUT` to protect against slow or incomplete headers;
- use `HTTP_READ_TIMEOUT` to limit the total time spent reading the full request;
- use both when you want protection for headers and an upper bound for the whole request.

For regular APIs, `HTTP_READ_HEADER_TIMEOUT` is often more important than a long `HTTP_READ_TIMEOUT`, because most API requests should send headers quickly even if the body is processed separately.

### Recommended values

For most HTTP APIs, a good default is:

- `2s–5s` for regular APIs and internal services;
- `5s–10s` for public APIs, mobile clients, or networks with higher latency;
- `10s+` only when there is a clear reason to tolerate very slow header delivery.

Longer values should be used carefully because they allow incomplete requests to hold server resources for longer.

## `HTTP_IDLE_TIMEOUT`

### Purpose

> `HTTP_IDLE_TIMEOUT` defines how long the server waits for the next request on an idle keep-alive connection.

After a response is completed, HTTP keep-alive lets the client reuse the same connection for another request. Reuse avoids repeated TCP and TLS setup, but an idle connection still consumes server resources.

### What it protects against

This timeout limits how long unused keep-alive connections remain open.

Examples:
- a client sends one request and leaves the connection open;
- a mobile client loses connectivity without closing the connection cleanly;
- a proxy keeps backend connections open but does not reuse them;
- many inactive clients consume available connection capacity.

### Behavior

When a keep-alive connection has no active request for the configured duration, the server closes it. Active requests are controlled by the read and write timeouts instead.

Go's `net/http.Server` falls back to `ReadTimeout` when `IdleTimeout` is zero. Configuring `HTTP_IDLE_TIMEOUT` explicitly keeps idle connection policy independent from request-reading policy.

### Recommended values

For most HTTP APIs, a practical range is:
- `30s–60s` for internal services or environments with many short-lived clients;
- `60s–120s` for typical public APIs;
- longer values only when reducing reconnect overhead is more important than releasing idle connections quickly.

Coordinate this value with reverse proxy and load balancer idle timeouts. Large differences can cause one side to close a connection while the other side still expects to reuse it.

## `HTTP_WRITE_TIMEOUT`

### Purpose

> `HTTP_WRITE_TIMEOUT` defines the maximum amount of time the server allows for writing an HTTP response to the client.

It protects the server from clients that receive responses too slowly or stop reading the response while keeping the connection open.

### What it protects against

This timeout is useful when the server has already processed the request and started sending a response, but the client reads the response too slowly or not at all.

Examples:
* the client has a very slow or unstable network connection;
* the client stops reading the response;
* the client downloads a large response extremely slowly;
* an attacker intentionally keeps response connections open to consume server resources.

### Why it matters

Writing a response also consumes server resources.

A slow-reading client may keep the connection open and force the server to hold resources such as:
* a socket;
* memory buffers;
* a goroutine or worker;
* connection capacity;
* response-writing state.

If many clients receive responses too slowly, the server may waste resources on already processed requests instead of serving normal users.

### Behavior

If the server cannot finish writing the response within the configured timeout, the connection is closed.

In simple terms: the server gives the client a reasonable amount of time to receive the response, but it does not allow the response-writing phase to hang forever.

### Important note

`HTTP_WRITE_TIMEOUT` is not a replacement for application-level request processing timeouts.

If the handler itself may run for too long, use a separate request context timeout or middleware-level timeout.

In Go HTTP servers, the write deadline can also affect long-running handlers, because the response must still be written within the configured write timeout.

### Recommended values

For most HTTP APIs, a good default is:
* `10s–30s` for regular JSON APIs and small responses;
* `30s–60s` for larger responses or slower public clients;
* `60s+` only when the service intentionally returns large responses or supports slow clients.

Be careful with very long values because they allow slow clients to hold server resources for longer.

For streaming responses, Server-Sent Events, WebSockets, or long-lived downloads, `HTTP_WRITE_TIMEOUT` may need special handling or a much higher value.

## `HTTP_SERVER_SHUTDOWN_TIMEOUT`

### Purpose

> `HTTP_SERVER_SHUTDOWN_TIMEOUT` defines the maximum amount of time the HTTP server is allowed to spend on graceful shutdown.

It controls how long the application waits for existing requests to finish before forcing the server to stop.

### What it protects against

This timeout is useful when the application receives a shutdown signal but still has active connections or in-flight requests.

Examples:
* the service is being stopped or restarted;
* the container is being terminated;
* a deployment is replacing the old application instance;
* the server has long-running requests that should be given time to finish;
* some clients keep connections open longer than expected.

### Why it matters

Without a shutdown timeout, the server may wait too long for active requests or connections to finish.

With a timeout, shutdown behavior becomes predictable:
* the server stops accepting new connections;
* existing requests get a chance to complete;
* idle connections are closed;
* the application does not hang forever during shutdown;
* deployments and restarts become safer and more controlled.

### Behavior

During graceful shutdown, the server stops accepting new requests and waits for active requests to complete.

If all active requests finish before the timeout expires, the server exits normally.

If the timeout expires first, `Shutdown` returns the context error. Active
connections are not forcibly interrupted by `net/http.Server.Shutdown`; the
application must decide whether forced termination is appropriate.

In simple terms: the server gives existing requests a final window to finish, but it does not wait forever.

### Important note

`HTTP_SERVER_SHUTDOWN_TIMEOUT` should usually be shorter than the platform-level termination grace period.

For example, in Kubernetes, the application should have enough time to:
* receive the termination signal;
* stop accepting new traffic;
* finish active requests;
* run graceful shutdown logic;
* exit before the container is forcefully killed.

### Recommended values

For most HTTP APIs, a good default is:
* `5s–10s` for small internal services with short requests;
* `15s–30s` for public APIs or services with moderate request duration;
* `30s–60s` for services that may have longer requests or heavier cleanup logic;
* `60s+` only when the service intentionally supports long-running requests or slow shutdown operations.

Very long shutdown timeouts should be used carefully because they can slow down deployments, restarts, and recovery from unhealthy instances.
