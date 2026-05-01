# Logging Skills Outline

## Goal

This document proposes how to split future AI agent skills for Go logging work in this repository style.

The split should reflect two valid operational modes:

- a simple baseline pattern
- a strict reference pattern

## Recommended Split

### Skill 1: Basic Logging Validation

Suggested purpose:

- implement or refactor logger config with `Validate()`
- keep bootstrap clean
- apply defaults in config types
- wire a concrete logger backend

Expected practices:

- `Validate()` on config structs
- config-owned defaults
- `main` stays thin
- logger backend constructors stay focused on backend behavior

When to use:

- small services
- internal tools
- low-complexity repos
- projects where decode-time strictness is not required

What the skill should enforce:

- no fake enums
- no reliance on `required` tags for runtime validation
- no backend-specific validation logic duplication

### Skill 2: Strict Logging Validation

Suggested purpose:

- implement or refactor logging config with `Validate() + UnmarshalText + decode hook`
- enforce enum-like safety during decode
- produce a reusable reference-grade config model

Expected practices:

- `Validate()` remains present
- enum-like types implement `UnmarshalText`
- `viper.Unmarshal(...)` uses a shared decode hook
- tests cover direct unmarshal and real decode-path behavior

When to use:

- platform or template repositories
- SDK-like internal foundations
- repos intended to teach or generate patterns
- systems where config correctness should fail as early as possible

What the skill should enforce:

- shared decode hook, not ad hoc hook wiring
- type-local validation in `UnmarshalText`
- object-level validation in `Validate()`
- no global mutable logger behavior when local alternatives exist

## Suggested Skill Selection Rule

Choose the basic skill by default.

Upgrade to the strict skill when at least one of the following is true:

- enum-like config types are reused across packages
- config is intended as a reference implementation
- early decode-time rejection is a project requirement
- agent-generated code should demonstrate best-practice strict parsing

## Suggested Deliverables for Each Skill

### Basic skill deliverables

- config structs
- `Validate()` methods
- defaults application
- logger bootstrap wiring
- tests for validation and defaults

### Strict skill deliverables

- everything from the basic skill
- `UnmarshalText` methods for enum-like types
- shared decode hook
- `viper.Unmarshal(...)` integration using that hook
- tests for direct `UnmarshalText`
- tests for invalid values failing during config decode

## Recommendation

The two-skill split is good.

Reason:

- it prevents pushing strict decode machinery into every simple service
- it preserves a high-quality reference pattern for cases where strictness is worth the extra moving parts
- it gives AI agents a clear escalation path instead of one oversized logging skill
