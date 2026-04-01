# Architecture Notes

This repository is planned as an API-first HR backend with two initial domains:

- `Core HR`
- `Leave & Attendance`

## Intended Top-Level Layout

```text
cmd/api
internal/platform
internal/modules/corehr
internal/modules/leave
migrations
docs/openapi
scripts
```

## Module Responsibilities

### `internal/platform`

Use for shared infrastructure only:

- config loading
- logging
- database bootstrap
- auth helpers
- middleware
- HTTP server setup

### `internal/modules/corehr`

Owns:

- employee records
- department and position relationships
- manager hierarchy
- assets and asset assignments

### `internal/modules/leave`

Owns:

- leave type definitions
- leave balances
- leave requests
- approval and rejection logic

## Cross-Module Rules

- `leave` may depend on manager relationships sourced from `corehr`.
- `corehr` should not depend on `leave`.
- Cross-module orchestration should stay explicit in services; do not hide it in ORM hooks.

## Operational Rules

- Keep API routes under `/api/v1`.
- Keep auth and role checks close to handlers and use cases.
- Keep migrations additive and reviewable.
