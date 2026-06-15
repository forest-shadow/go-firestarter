# Validation Strategy

This document captures the repository-level validation approach and the tradeoffs
around manual `Validate()` methods versus a general-purpose validator such as
`go-playground/validator/v10`.

For an implementation-oriented AI agent skill outline based on these practices,
see [`skills.md`](skills.md).

The goal is not to declare one universal rule. The goal is to keep validation
choices intentional and proportionate to the shape of the codebase.

## Current Default

For this microservice starter, configuration validation is currently treated as
runtime/domain validation rather than as a generic infrastructure concern.

Config values define the conditions under which the service can start and
operate correctly. Those rules are small, project-specific, and usually loaded
once at startup.

Current examples:

- `app.name` must be present
- `app.version` must be present
- `app.env` must be one of the supported environments
- `logger.level` must be one of the supported log levels
- `logger.format` must be one of the supported output formats
- default logger format depends on the selected application environment

These rules are easier to audit when they are expressed directly through
`Validate()`, `WithDefaults()`, `IsValid()`, and type-level parsing methods such
as `UnmarshalText()`.

## Domain Validation

Treat validation as domain validation when the rule expresses what is valid for
this application or this model.

Good signs:

- the rule uses project vocabulary, such as `AppEnv`, `LogLevel`, or `LogFormat`
- the rule protects runtime correctness
- the rule depends on project policy, defaults, or allowed combinations
- the rule is easier to understand as code than as a struct tag
- the model is validated in a few deliberate places instead of many repeated
  input contracts

Examples:

- `local` and `development` default to console logging
- other environments default to JSON logging
- a config object must be complete before the service starts
- enum-like values should reject unsupported project-specific values
- a later domain model cannot transition into an invalid state

Preferred tools:

- explicit `Validate()` methods
- explicit `WithDefaults()` methods
- enum-like `IsValid()` methods
- `UnmarshalText()` for typed scalar parsing

## Infrastructure Validation

Treat validation as infrastructure validation when the rule is mostly a generic
input constraint repeated across many structs.

Good signs:

- many request, form, query, or webhook DTOs repeat similar checks
- validation rules are generic constraints such as required, email, URL, UUID,
  min length, max length, or numeric ranges
- validation errors need a centralized API response format
- the project benefits from one shared validation mechanism across many input
  contracts
- the same validation engine is used by multiple packages or transport layers

Examples:

- HTTP request bodies
- query parameter structs
- CLI input structs
- webhook payloads
- admin form payloads
- public or internal API DTOs

In those cases, a library such as `go-playground/validator/v10` can be a good
fit because it turns repeated mechanical validation into shared infrastructure.

## Valid Approaches

The following approaches are all valid when chosen intentionally.

### Manual validation for config and domain models

Use manual validation when the rules are small, project-specific, and tied to
runtime correctness.

This is the current default for the starter config.

Benefits:

- clear project-specific intent
- no extra dependency
- precise errors
- easy to combine defaults, enum-like types, and cross-field checks

Tradeoffs:

- repeated generic checks can become boilerplate if the number of input structs
  grows
- error formatting is owned by the project

### Project-wide validator usage

Use a general-purpose validator everywhere when the project intentionally
standardizes on one validation mechanism.

This can include config, API DTOs, form payloads, and other input models.

Benefits:

- one validation style
- one error extraction path
- useful for larger teams and large input surfaces
- strong fit for many generic constraints

Tradeoffs:

- domain-specific rules can become less visible when hidden behind tags or
  custom validators
- config policy may be harder to read than explicit code
- custom error mapping and custom validators become part of the project
  infrastructure

### Hybrid validation

Use a hybrid model when both kinds of rules exist.

Recommended split:

- use a general-purpose validator for repeated external input contracts
- use manual `Validate()` methods for runtime config and domain invariants
- keep complex business or lifecycle rules in explicit code

This is often the most pragmatic option once the service grows an API layer.

## Decision Point

Consider introducing a general-purpose validator when validation becomes a
repeated infrastructure task rather than a small set of model-specific rules.

Practical triggers:

- many DTOs repeat the same required, range, length, URL, email, UUID, or enum
  checks
- multiple packages need the same validation flow
- API responses need centralized validation error formatting
- validation tags would remove meaningful boilerplate without hiding important
  project policy
- custom validators are shared by several input contracts

Stay with manual validation when:

- the config is small and loaded once at startup
- the rules define runtime invariants for this service
- the rules use project-specific vocabulary
- explicit code is clearer than generic tags
- there is no repeated validation surface yet

## Repository Guidance

For this microservice starter:

- keep config validation manual by default
- keep enum-like parsing on the enum-like types
- keep `Validate()` as the final mandatory validation layer
- consider `go-playground/validator/v10` later for API/request validation if the
  project gains enough repeated input contracts
- avoid adding a validation framework only to validate a small startup config
