# `internal`

## Purpose

`internal` contains repository-private application code.

It is the place for packages that belong to this service or starter layout but
should not become a reusable public API for external importers.

It currently owns:

- `internal/app` for application bootstrap and dependency wiring

## Package boundaries

Use `internal/*` for code that coordinates reusable packages or depends on this
repository's runtime layout.

Good candidates include:

- startup composition
- integration glue
- application-specific file naming and environment fallback rules
- concrete dependency selection

Keep reusable vocabulary, config models, and implementation packages in `pkg/*`
when they are intended to be shared by multiple callers or imported outside the
application bootstrap path.

## Import behavior

Go's `internal` import rules prevent packages outside the parent tree from
importing code below this directory. That makes it a useful boundary for
implementation details that should remain private to the repository.
