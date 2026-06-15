# Validation Skills Outline

## Goal

This document proposes a practical AI agent skill for implementing validation in
this repository style.

The skill should help an agent choose and implement the right validation mode
instead of applying a validation library by default.

It should cover:

- repository-local config validation practices
- domain/runtime invariants
- repeated input-contract validation
- hybrid validation for growing services
- decode-hook based parsing and validation for typed scalar values

## Skill: Practical Go Validation

Suggested purpose:

- implement or refactor validation for config, domain models, or input DTOs
- choose between manual `Validate()` methods, `go-playground/validator/v10`, or a hybrid approach
- keep runtime config validation explicit and easy to audit
- introduce shared validation infrastructure only when repeated input contracts justify it
- preserve typed scalar parsing through `UnmarshalText()` and decode hooks where strict decode behavior is useful

## Repository Default

For this starter, the default is manual validation for runtime config.

Expected practices:

- config structs expose explicit `Validate()` methods
- config defaults live in explicit methods such as `WithDefaults()`
- enum-like values expose `IsValid()`
- reusable enum-like or scalar types implement `UnmarshalText()` when decode-time validation is useful
- `viper`/`mapstructure` integration uses a shared decode hook instead of ad hoc hook wiring
- the startup path runs final object-level validation after defaults are applied

Current examples:

- `pkg/config.App.Validate()`
- `pkg/logger.Config.Validate()`
- `pkg/logger.Config.WithDefaults(...)`
- `env.AppEnv.UnmarshalText(...)`
- `logger.LogLevel.UnmarshalText(...)`
- `logger.LogFormat.UnmarshalText(...)`
- `pkg/config.DecodeHook()`

## Mode A: Manual Domain Validation

Use this mode when validation expresses runtime correctness or model-specific
invariants.

Good fit:

- startup config
- project-specific enum-like values
- small domain models
- rules involving defaults or allowed field combinations
- invariants that should be readable as ordinary Go code

Expected implementation:

- add or update `Validate()` on the model
- keep error messages specific to the config or domain path
- add `WithDefaults()` when defaults are part of the model contract
- add `IsValid()` for enum-like values
- add table tests for valid and invalid cases

Avoid:

- hiding project policy inside generic tags
- introducing a validation framework for a small startup config
- relying on struct tags such as `required` unless a validator actually enforces them

## Mode B: Validator-based Infrastructure Validation

Use this mode when validation is a repeated infrastructure concern across many
input contracts.

Good fit:

- HTTP request DTOs
- query parameter structs
- form payloads
- webhook payloads
- CLI input structs
- many structs repeating generic constraints

Expected implementation:

- introduce a single validator setup point
- keep validator usage behind a small application or transport-layer helper
- use tags for generic constraints such as `required`, `email`, `url`, `uuid`,
  `min`, `max`, `gte`, `lte`, and `oneof`
- centralize validation error formatting for the transport layer
- register custom validators only when the rule is reused by several input contracts

Avoid:

- using validator tags for complex domain policy that is clearer in code
- scattering validator instances across packages
- mixing API error formatting into domain models
- replacing typed enum parsing with untyped strings only to satisfy tag-based validation

## Mode C: Hybrid Validation

Use this mode when the service has both runtime/domain invariants and repeated
external input contracts.

Recommended split:

- manual `Validate()` for config and domain invariants
- `go-playground/validator/v10` for repeated DTO/form/request validation
- explicit service/domain code for business rules and lifecycle transitions
- decode hooks for strict config parsing when typed scalar values should reject invalid input during decode

Example flow:

1. Decode config using the shared decode hook.
2. Apply defaults.
3. Run manual config `Validate()`.
4. Validate API DTOs through the shared validator helper.
5. Enforce business invariants in service/domain code.

This keeps each validation layer responsible for the kind of rule it handles
best.

## Decode-hook Pattern

Use the decode-hook pattern when enum-like or scalar types should validate
themselves during config loading.

Expected implementation:

- define a small typed scalar or enum-like type
- implement `IsValid()` for explicit checks
- implement `UnmarshalText()` for decode-time parsing
- expose one shared decode hook from the config loading package
- pass that hook into `viper.Unmarshal(...)`
- keep final `Validate()` methods even when decode-time validation exists

Why final `Validate()` still stays:

- not every construction path uses the decode hook
- object-level invariants may involve multiple fields
- defaults are usually applied after decode
- final validation makes startup correctness explicit

Avoid:

- relying on `UnmarshalText()` without wiring `mapstructure.TextUnmarshallerHookFunc()`
- duplicating decode hook composition at each call site
- moving cross-field rules into scalar parsers

## Selection Rule

Start with manual validation.

Move to validator-based infrastructure only when at least one of the following
is true:

- many input structs repeat the same generic constraints
- multiple packages need the same validation flow
- API validation errors need one centralized response format
- validator tags remove meaningful boilerplate without hiding project policy
- custom validators are reused across several input contracts

Use the hybrid mode when the project has both small runtime config invariants and
a larger API/input surface.

## Suggested Deliverables

For manual config/domain validation:

- `Validate()` methods
- defaults methods when needed
- enum-like `IsValid()` methods
- optional `UnmarshalText()` for strict decode behavior
- table tests for direct validation
- decode-path tests when a decode hook is used

For validator-based input validation:

- validator setup helper
- DTO tags for generic constraints
- centralized validation error extraction
- transport-level tests for invalid payloads
- custom validators only when rules are shared

For hybrid validation:

- clear package boundary between domain/config validation and DTO validation
- documentation explaining which layer owns which rules
- tests that exercise both generic input validation and domain/runtime validation

## Recommendation

For this repository, keep runtime config validation manual and explicit.

Introduce `go-playground/validator/v10` later only if the service grows enough
API or input DTOs to make validation a repeated infrastructure concern.

Keep the decode-hook approach for strict typed config parsing when enum-like
values should fail during unmarshal.
