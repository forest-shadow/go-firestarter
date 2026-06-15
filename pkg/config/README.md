# `pkg/config`

## Purpose

`pkg/config` is the home of the repository's typed configuration model.

It currently owns:

- config structs such as `App`
- struct-level validation via `Validate()`
- the shared decode hook used by the current `viper`/`mapstructure` pipeline

It intentionally uses shared starter vocabulary from neighboring packages:

- `pkg/env` provides `AppEnv`

This package is optimized for reference clarity and a stable starter config contract rather than for the strictest possible layer separation.

## Why the decode hook is here

In a more heavily layered architecture, `DecodeHook()` could live in an internal integration package because it belongs to config loading, not to the pure config model.

For this repository, it is kept here on purpose.

Reason:

- this repository is a starter/reference pack
- keeping config validation and decode integration in one place makes them easier to study, reuse, and simplify
- if strict decode-time validation later proves unnecessary, the unmarshalling layer can be removed more easily when it is centralized

Enum-like scalar types still own their own `UnmarshalText()` behavior in their
domain packages. `pkg/config` provides the shared decode hook that enables that
behavior in the current `viper`/`mapstructure` pipeline.

## When to consider moving decode integration into `internal`

Consider moving `DecodeHook()` out of `pkg/config` and into an internal package such as `internal/appconfig` when:

- `pkg/config` starts to feel too infrastructure-heavy
- config loading needs diverging decode-hook behavior
- the codebase gains multiple config-loading paths
- `pkg/config` should become a cleaner model-only package
- the repository prioritizes stricter layer isolation over reference convenience

If that happens, keep the following in `pkg/config`:

- config types
- `Validate()`

Move only the framework-specific decode integration.
