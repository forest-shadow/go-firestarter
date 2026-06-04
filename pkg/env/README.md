# `pkg/env`

## Purpose

`pkg/env` defines the repository's shared application environment vocabulary.

It currently owns:

- the `AppEnv` type
- supported environment constants such as `local`, `development`, `staging`, and `production`
- validation helpers for checking environment values
- text unmarshalling for config loaders that use `encoding.TextUnmarshaler`

The package keeps environment names centralized so config, logging, telemetry, and application setup code can share the same stable values.

## Environment values

The supported values are:

- `local`
- `development`
- `staging`
- `production`

Use `IsValid()` when accepting an `AppEnv` value from outside the type system, such as configuration files, environment variables, or tests that construct values directly.

Use `IsLocal()` for behavior that should only apply to local development.

## Config loading

`AppEnv` implements `UnmarshalText()` so configuration loaders can decode text values into the typed environment model.

In this repository, `pkg/config.DecodeHook()` enables that behavior in the current `viper`/`mapstructure` pipeline.
