# `pkg/logger`

## Purpose

`pkg/logger` defines the repository's shared logging configuration vocabulary
and runtime logging contract.

It currently owns:

- the `LogLevel` type
- supported log levels such as `debug`, `info`, `warn`, and `error`
- the `LogFormat` type
- supported output formats such as `json` and `console`
- the `Field` type and `F()` helper for structured log fields
- the `Logger` interface used by application code
- validation helpers for checking logging config values
- text unmarshalling for config loaders that use `encoding.TextUnmarshaler`

The package keeps logging config values and the runtime logging facade
centralized so application code can depend on a stable contract while backend
implementations remain replaceable.

## Log levels

The supported log levels are:

- `debug`
- `info`
- `warn`
- `error`

Use `LogLevel.IsValid()` when accepting a level from outside the type system, such as configuration files, environment variables, or tests that construct values directly.

## Log formats

The supported output formats are:

- `json`
- `console`

Use `LogFormat.IsValid()` when accepting a format from configuration or other external input.

## Config loading

`LogLevel` and `LogFormat` implement `UnmarshalText()` so configuration loaders can decode text values into typed logging settings.

In this repository, `pkg/config.DecodeHook()` enables that behavior in the current `viper`/`mapstructure` pipeline.

## Runtime API

Application code should depend on the `Logger` interface instead of concrete
backend types.

Use `F()` to attach structured fields:

```go
log.Info("application started", logger.F("env", env))
```

The interface intentionally exposes only the common operations this starter
needs:

- level-based logging through `Debug`, `Info`, `Warn`, and `Error`
- structured fields through `Field`
- child loggers through `With`
- backend flushing through `Sync`

### `Sync`

`Sync()` is part of the shared logger contract so application entrypoints can
flush backend-owned buffers before shutdown without depending on a concrete
logger implementation.

This matters most for backends such as `zap`, where log output may pass through
buffered or syncable writers. Calling `Sync()` gives the backend a final chance
to write pending log entries before the process exits. Some implementations may
not need any flushing; for example, the current `zerolog` adapter returns `nil`.

When logging to standard streams, some backends or platforms may return a
non-actionable sync error even after logs were written. Application code should
still call `Sync()` at shutdown, but it can choose whether to ignore, filter, or
handle that error based on the selected output and runtime requirements.

## Logger implementations

Concrete logger constructors live in subpackages:

- `pkg/logger/zaplogger`
- `pkg/logger/zerologger`

Those packages consume `pkg/logger` values through `pkg/config.Logger`, apply defaults, validate the selected level and format, and build the concrete logger.
Both constructors return the shared `Logger` interface.
