# `internal/app`

## Purpose

`internal/app` is the repository's application bootstrap package.

It currently owns:

- loading runtime configuration from the selected `env.<APP_ENV>.yml` file
- applying environment variable overrides through the current `viper` pipeline
- validating the assembled application configuration
- applying logger defaults based on the selected app environment
- applying HTTP server defaults
- constructing the concrete application logger

This package is intentionally internal because it wires repository-specific
startup behavior rather than defining reusable library contracts.

## Package boundaries

`internal/app` coordinates reusable packages instead of owning their vocabulary.

It depends on:

- `pkg/config` for app config structures, validation, and decode hooks
- `pkg/env` for supported application environment values
- `pkg/httpserver` for HTTP server configuration and lifecycle behavior
- `pkg/logger` for logger configuration, vocabulary, and the shared logger contract
- `pkg/logger/zaplogger` for the current concrete logger backend

Keep reusable config models, enum-like values, and logger implementation details
in `pkg/*`. Keep application-specific composition, file naming, environment
fallbacks, and startup wiring here.

## Config loading

`GetConfig()` reads `APP_ENV` to choose the config file name. When `APP_ENV` is
not set, it falls back to `local` and loads `env.local.yml`.

The loader:

- reads YAML config from the repository root
- enables `APP_`-prefixed environment variable overrides
- decodes scalar config values through `pkg/config.DecodeHook()`
- validates app, logger, and HTTP server config before returning

For example, `http_server.read_timeout: 20s` can be overridden with
`APP_HTTP_SERVER_READ_TIMEOUT=20s`, and `http_server.read_header_timeout: 10s`
can be overridden with `APP_HTTP_SERVER_READ_HEADER_TIMEOUT=10s`. The same
mapping applies to `http_server.idle_timeout` and
`APP_HTTP_SERVER_IDLE_TIMEOUT`.

## Logger setup

`NewLogger()` builds the current application logger from the validated config
and returns the shared `pkg/logger.Logger` contract.

The selected implementation is `pkg/logger/zaplogger`. Other logger
implementations can remain reusable in `pkg/logger/*`, while the application
chooses which one to wire here.

This package intentionally keeps the backend choice explicit instead of adding a
separate unified logger-constructor config too early. The small mapping from the
application `Config` into `zaplogger.Config` makes `internal/app/logger.go` the
single place that knows which concrete logger is used. The rest of the
application should continue to depend only on the shared `pkg/logger.Logger`
API.
