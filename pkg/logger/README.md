# `pkg/logger`

## Purpose

`pkg/logger` defines the repository's shared logging configuration vocabulary.

It currently owns:

- the `LogLevel` type
- supported log levels such as `debug`, `info`, `warn`, and `error`
- the `LogFormat` type
- supported output formats such as `json` and `console`
- validation helpers for checking logging config values
- text unmarshalling for config loaders that use `encoding.TextUnmarshaler`

The package keeps logging config values centralized so application config and logger implementations can share the same stable contract.

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

## Logger implementations

Concrete logger constructors live in subpackages:

- `pkg/logger/zaplogger`
- `pkg/logger/zerologger`

Those packages consume `pkg/logger` values through `pkg/config.Logger`, apply defaults, validate the selected level and format, and build the concrete logger.
