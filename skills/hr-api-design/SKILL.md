---
name: hr-api-design
description: Use this skill when designing, reviewing, or extending the REST API for this HR backend. It captures the project’s API-first direction, resource groups, naming conventions, and guidance for stable request and response design without over-promising unfinished details.
---

# HR API Design

Use this skill when adding or reviewing HTTP endpoints for this HR project.

## API Direction

- The backend is `REST JSON`.
- Base path is `/api/v1`.
- The API should support future web admin and Flutter clients.

## Primary Resource Groups

- `auth`
- `employees`
- `organization`
- `assets`
- `leave-types`
- `leave-balances`
- `leave-requests`

## Route Design Rules

- Prefer resource-oriented naming.
- Use plural resource names for collections.
- Keep action-style routes for domain actions that are not simple CRUD, such as approvals.
- Keep URLs stable and move business branching into services.

## Request and Response Guidance

- Use JSON consistently.
- Keep external DTOs clear and purpose-built.
- Do not expose internal persistence quirks in the public API.
- Do not freeze exact payload shapes in docs until the implementation is real.

## Error Handling

- Return predictable error shapes.
- Separate validation, authentication, authorization, and domain-state failures.
- Keep internal error details out of client-facing messages.

## Authorization Guidance

- Check authentication at the middleware boundary.
- Check role and ownership constraints in handlers and services.
- For leave approval, verify both role and reporting-line authority.

## Workflow-Oriented Endpoints

Prefer explicit action routes for state transitions:

- `POST /leave-requests/:id/approve`
- `POST /leave-requests/:id/reject`
- `POST /assets/assign`
- `POST /assets/return`

## Documentation Rules

- Keep OpenAPI aligned with real handlers.
- Do not document speculative routes as committed contracts.

## When To Read References

- Read [API Conventions](references/api-conventions.md) when defining new endpoints or response patterns.
