# Logging and Config Validation

## Purpose

This document captures the current logging architecture and the validation practices used around configuration loading.

The goal is to provide a stable reference for:

- application logging bootstrap
- logger backend implementation patterns
- the shared logging config contract
- the shared runtime logger contract
- config validation with `Validate()`
- strict enum-like config parsing with `UnmarshalText`
- `viper`/`mapstructure` decode hooks
- future AI agent skills for logging-related work

## Current Architecture

### Bootstrap flow

The application currently follows this startup sequence:

1. `internal/app.GetConfig()` loads config from file and environment.
2. `viper.Unmarshal(...)` decodes into typed config structs.
3. `config.DecodeHook()` enables `encoding.TextUnmarshaler`-based parsing.
4. `App.Validate()` validates the app section.
5. `logger.Config.WithDefaults(...)` applies logger defaults.
6. `logger.Config.Validate()` validates the logger section.
7. `internal/app.NewLogger(...)` constructs the selected logger backend and returns the shared logger contract.

This keeps `cmd/main.go` thin and prevents logging setup details from leaking into the entrypoint.

### Package responsibilities

- `internal/app`
  Owns application bootstrap and config loading.
- `pkg/config`
  Owns shared config-loading support and app-level config structures.
- `pkg/env`
  Owns the starter's environment vocabulary.
- `pkg/logger`
  Owns the starter's shared logging vocabulary, logger config, runtime logger contract, and default logging policy.
- `pkg/logger/zaplogger`
  Owns the concrete `zap` backend adapter.
- `pkg/logger/zerologger`
  Owns the alternative `zerolog` backend adapter.

This split is intentionally pragmatic. The starter defines an opinionated config contract and a small runtime logger interface, and the concrete logger backends implement those contracts.

### Layering model

The intended dependency direction is:

```text
pkg/env
  <- pkg/config
  <- pkg/logger
  <- pkg/logger/zaplogger, pkg/logger/zerologger
  <- internal/app
```

`pkg/config`, `pkg/env`, and `pkg/logger` are not treated as fully independent libraries. Together, they define the starter's reusable defaults. More specific applications can replace or adapt this contract when they need different config structure, environment vocabulary, or logger policy.

## Logging Design Principles

### 1. Keep the entrypoint small

The `main` package should not manually assemble logger internals. It should call a bootstrap-level constructor such as `app.NewLogger(...)`.

Why:

- reduces startup noise
- centralizes operational wiring
- makes config and logger initialization easier to test

### 2. Keep config rules and logging policy explicit

Config validation should live next to the config types, not inside logger backends.

Examples:

- `logger.Config.WithDefaults(appEnv)`
- `logger.Config.Validate()`
- `App.Validate()`
- `logger.DefaultLogFormat(appEnv)`

Why:

- one source of truth
- avoids duplicating config rules across `zap` and `zerolog`
- makes behavior predictable across backends

`logger.DefaultLogFormat(appEnv)` intentionally lives in `pkg/logger`: choosing `console` for local/development and `json` elsewhere is logging policy, while `logger.Config.WithDefaults(...)` is the place where that policy is applied to the loaded config.

### 3. Let backend packages own backend behavior

Logger-specific packages should focus on:

- level conversion
- encoder/output selection
- backend-specific formatting
- adapting backend APIs to the shared `logger.Logger` interface

They should not be the primary home for business rules about config validity.

### 4. Keep application logging behind the shared contract

Application packages should accept `logger.Logger` when they need to log.

Why:

- avoids package-global logger state
- keeps `zap` and `zerolog` replaceable
- prevents backend method names from leaking into application code
- makes logger dependencies explicit in constructors

### 5. Prefer explicit operational behavior

Good examples from the current implementation:

- use `stderr` for log output
- switch default format by environment
- colorize console output only when running in a TTY
- add stable app metadata fields like `app_name` and `app_version`

These are concrete, operationally useful defaults.

## Validation Strategy

For the repository-level validation tradeoffs, including when to prefer manual
`Validate()` methods and when a general-purpose validator becomes useful, see
[`docs/validation`](../validation/README.md).

### Layer 1: struct-level validation with `Validate()`

This is the baseline and most important practice.

Use `Validate()` when:

- a field must not be empty
- a value must belong to an allowed set
- multiple fields must be checked together
- the config is internal and loaded once at startup

Examples:

- `App.Validate()`
- `logger.Config.Validate()`

Typical rules:

- `app.name` is required
- `app.version` is required
- `app.env` must be one of the known environments
- `logger.level` must be one of the supported levels
- `logger.format` must be one of the supported formats

### Layer 2: type-level validation with `UnmarshalText`

Use `UnmarshalText` for enum-like or scalar domain types that should reject invalid values during decode.

Examples:

- `env.AppEnv`
- `logger.LogLevel`
- `logger.LogFormat`

This layer is useful when a type:

- is reused in more than one config struct
- should protect itself independently of any parent struct
- benefits from fail-fast parsing

### Why both layers can coexist

They solve different problems:

- `UnmarshalText` validates one value in isolation during decode
- `Validate()` validates a completed config object after decode

`Validate()` should remain even when `UnmarshalText` exists, because:

- object-level invariants still need a home
- decode hooks may not be used everywhere
- validation remains explicit and easy to audit

## `UnmarshalText` and Decode Hooks

### Important nuance

Implementing `UnmarshalText` alone is not enough for `viper.Unmarshal(...)`.

`viper` uses `mapstructure`, and `TextUnmarshaler` support must be explicitly enabled via a decode hook.

### Project pattern

The project exposes a shared decode hook:

```go
func DecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToWeakSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
	)
}
```

And uses it during config decode:

```go
if err := v.Unmarshal(&cfg, viper.DecodeHook(config.DecodeHook())); err != nil {
	return nil, fmt.Errorf("decode config to struct: %w", err)
}
```

This is the strict path: invalid enum-like values fail during unmarshal instead of only during later validation.

### Why `DecodeHook()` currently stays in `pkg/config`

In a stricter layered architecture, `DecodeHook()` could live in an integration package such as `internal/appconfig`, because it is tied to the `viper`/`mapstructure` loading pipeline rather than to the pure domain model.

This repository intentionally keeps `DecodeHook()` in `pkg/config` for now.

Reason:

- this repository is primarily a reusable starter/reference pack
- keeping defaults, `Validate()`, `UnmarshalText`, and decode-hook wiring close together makes the pattern easier to study and copy
- if strict decode-time validation later proves unnecessary, the unmarshalling layer is easier to remove when all config-validation practices are documented and implemented in one place

This is a convenience-oriented reference decision, not a universal rule.

### When to consider moving decode integration into `internal`

Revisit the placement of `DecodeHook()` if one or more of the following becomes true:

- `pkg/config` starts accumulating framework-specific loading concerns
- the project introduces more than one config-loading pipeline
- different applications need different decode-hook compositions
- the package is becoming harder to understand because model rules and loader integration are growing together
- the repository shifts from reference convenience toward stricter layer isolation

If those conditions appear, a good next step is to move `mapstructure`/`viper` decode integration into an internal package while keeping `Validate()` and `UnmarshalText` on the config types themselves.

## Choosing Between Two Validation Modes

### Mode A: `Validate()` only

Use this when:

- the project is small or internal
- config is loaded once at startup
- types are local and not widely reused
- simplicity matters more than strict decode behavior

Benefits:

- simple mental model
- low framework coupling
- easy to read and maintain

Tradeoffs:

- invalid enum-like values are caught later
- type-level invariants are less reusable

### Mode B: `Validate()` + `UnmarshalText` + decode hook

Use this when:

- config quality matters operationally
- enum-like types are reused
- the project wants strict decode behavior
- the codebase is intended as a reusable reference

Benefits:

- invalid values fail earlier
- enum-like types become self-validating
- better reference architecture for future generators and agents

Tradeoffs:

- slightly more moving parts
- `viper`/`mapstructure` integration must be maintained deliberately

### Recommended default for this repository

Keep the current strict mode:

- `Validate()` stays mandatory
- `UnmarshalText` stays on enum-like types
- `config.DecodeHook()` stays the canonical decode path

This gives the repository reference value without making the runtime model overly complex.

## Backend Implementation Practices

### `zaplogger`

Current good practices:

- `New(Config)` is small and explicit
- `stderr` is wrapped with `zapcore.Lock(...)`
- JSON and console output are clearly separated
- console color is enabled only for TTY output
- backend adds stable app metadata fields

### `zerologger`

Current good practices:

- reuses the same validated config model
- remains an alternative backend, not a forked config design
- avoids mutating global `zerolog.TimeFieldFormat`

Rule:

Alternative backends should share config semantics even if their output formatting differs.

### Compile-time interface checks

Concrete logger backends should declare a compile-time interface check near the
implementation type:

```go
var _ logger.Logger = Logger{}
```

This does not create a runtime dependency or register the implementation
anywhere. It asks the compiler to verify that the concrete backend type still
satisfies the shared `logger.Logger` contract.

Use this pattern when a package intentionally implements a public interface but
the relationship is otherwise only visible through constructors or tests. If the
interface changes, or if the backend method signatures drift, the package fails
to compile at the implementation boundary instead of failing later at a call
site.

## Anti-Patterns to Avoid

- putting config validation only inside backend constructors
- relying on struct tags like `required:"true"` as if `viper.Unmarshal` enforces them
- assuming typed string aliases are true enums
- adding color codes to non-TTY output
- mutating package-global logger state when a local mechanism is possible
- leaking backend-specific APIs such as `Infow` or `zerolog.Event` into application packages
- expanding `logger.Logger` into a broad clone of a backend API

## Testing Guidance

At minimum, test:

- default application of logger config
- `Validate()` success and failure paths
- `UnmarshalText` success and failure paths
- `viper.Unmarshal(...)` with the shared decode hook

The most important integration test is not the method itself, but the real decode path using `viper.DecodeHook(config.DecodeHook())`.

## Guidance for Future Refactors

The codebase now uses a backend-independent logger facade.

When evolving it:

- keep `internal/app.NewLogger(...)` returning `logger.Logger`
- keep the interface focused on capabilities used by the application
- keep logger config ownership in `pkg/logger`
- do not duplicate validation logic per backend
- add adapter behavior in backend packages, not in callers

## Summary

The current recommended standard for this repository is:

- thin bootstrap in `internal/app`
- config-owned defaults and validation
- enum-like types guarded by `UnmarshalText`
- explicit decode hook integration for `viper`
- backend-specific logic isolated in backend packages
- `Validate()` retained as the final and mandatory validation layer
