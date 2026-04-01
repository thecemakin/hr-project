---
name: hr-backend-implementation
description: Use this skill when implementing or refactoring the Go HR backend in this repository. It defines the expected modular monolith structure, package boundaries, persistence approach with PostgreSQL and GORM, and the default implementation style for Core HR and Leave modules.
---

# HR Backend Implementation

Use this skill when working on the Go backend for this HR project.

## Goals

- Keep the codebase as a modular monolith.
- Preserve clean boundaries between `platform`, `corehr`, and `leave`.
- Favor readable business logic over generic abstractions.
- Keep the API-first shape stable for future web and Flutter clients.

## Architecture Rules

- Put startup wiring in `cmd/api`.
- Put shared infrastructure in `internal/platform`.
- Put business logic inside the owning module under `internal/modules`.
- Do not place HR-specific rules in generic helpers or middleware.
- Prefer explicit dependencies over package-level globals.

## Default Package Shape

For each business module, prefer a shape like:

- `handler` or transport layer for HTTP mapping
- `service` for business rules
- `repository` for persistence
- `model` or entity definitions close to the module

Adjust naming to the repo style once the first implementation is established, but keep responsibilities separate.

## Persistence Guidance

- Use `PostgreSQL` as the source of truth.
- Use `GORM` for straightforward CRUD and relationship mapping.
- Drop to explicit SQL in repositories when reporting or filtering logic becomes awkward in GORM.
- Keep transaction boundaries in the service layer when the use case updates multiple records.

## Domain Boundaries

- `corehr` owns employees, organization relationships, and asset lifecycle.
- `leave` owns leave types, balances, requests, and approval flow.
- `users` and auth are platform-adjacent, but their role mapping must support HR use cases.

Read [Architecture Notes](references/architecture.md) before making large structural changes.

## Implementation Style

- Prefer small services with explicit method names.
- Return domain-friendly errors and map them to HTTP responses at the handler layer.
- Avoid premature interfaces unless there is a real second implementation or a test seam that benefits from it.
- Keep request/response DTOs separate from database models when external API shape starts to diverge.

## Validation Priorities

- Protect authorization boundaries first.
- Validate manager relationships before approval actions.
- Validate leave balance before final approval.
- Preserve asset assignment history; do not overwrite the past state destructively.

## Testing Expectations

- Add unit tests for domain rules.
- Add integration tests for module flows that touch the database.
- Prefer scenario tests for approval, balance deduction, and access control.

## When To Read References

- Read [Architecture Notes](references/architecture.md) when adding new packages, modules, or cross-module flows.
