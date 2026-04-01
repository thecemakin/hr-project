# API Conventions

## Base Path

All v1 routes should live under:

```text
/api/v1
```

## Example Route Groups

```text
POST   /api/v1/auth/login

GET    /api/v1/employees
POST   /api/v1/employees
GET    /api/v1/employees/:id
PATCH  /api/v1/employees/:id

GET    /api/v1/organization/tree

GET    /api/v1/assets
POST   /api/v1/assets
POST   /api/v1/assets/assign
POST   /api/v1/assets/return

GET    /api/v1/leave-types
GET    /api/v1/leave-balances/:employeeId
POST   /api/v1/leave-requests
GET    /api/v1/leave-requests/me
GET    /api/v1/leave-requests/pending-approvals
POST   /api/v1/leave-requests/:id/approve
POST   /api/v1/leave-requests/:id/reject
```

## Status Guidance

Favor a small set of clear domain statuses. Examples:

- leave request: `pending`, `approved`, `rejected`
- asset assignment: `assigned`, `returned`
- employee: `active`, `inactive`

## DTO Guidance

- Separate transport DTOs from persistence models when needed.
- Keep date and timestamp formats consistent across modules.
- Prefer explicit field names over overloaded generic payloads.
