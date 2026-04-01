# HR Backend

`HR Backend` is a Go-based human resources platform focused on building a solid backend foundation for two high-value modules first: `Core HR` and `Leave & Attendance`. The initial release is designed as an API-first, single-company system that gives teams a reliable digital employee record, a clear organizational structure, asset assignment tracking, and a streamlined leave approval flow that can later power both a web admin interface and a Flutter mobile app.

## Why This Project Exists

Many HR processes still depend on spreadsheets, shared folders, manual approvals, and scattered communication. That makes it hard to answer simple questions reliably:

- Who reports to whom?
- Which laptop is assigned to which employee?
- How many annual leave days does someone have left?
- Which leave requests are waiting for approval?

This project aims to bring those workflows into a single backend system with clear ownership, clean domain boundaries, and room to grow into a broader HR platform over time.

## v1 Scope

The first version intentionally focuses on a narrow but valuable slice of the product:

- Single-company deployment
- API-first backend
- REST JSON interface
- Email/password authentication with JWT
- Four application roles: `Admin`, `HR`, `Manager`, `Employee`
- Single-level leave approval based on direct manager relationship

### Included in v1

- Employee master records
- Digital personnel file foundation
- Organization hierarchy lookup
- Asset assignment and return tracking
- Leave type definitions
- Leave balance tracking
- Leave request submission
- Manager approval or rejection flow

### Explicitly Deferred

- Payroll
- Recruitment / ATS
- Performance management
- Multi-company SaaS isolation
- Advanced accrual rules by region or tenure
- Multi-step approval engines
- Notification delivery infrastructure
- Final web admin panel
- Final Flutter mobile client

## Core Modules

### 1. Core HR

`Core HR` is the system of record for employee-related information. In product terms, it acts as the digital employee card and the structural backbone for the rest of the platform.

It is expected to cover:

- Personal identity and contact information
- Address information
- Emergency contacts
- Bank account details
- Employment metadata such as department, position, status, and manager
- Organization hierarchy relationships
- Assigned company assets such as laptops or phones

From a backend perspective, this module owns the employee lifecycle, manager relationships, and asset assignment state.

### 2. Leave & Attendance

`Leave & Attendance` addresses one of the most automation-friendly HR workflows: leave requests and approval tracking.

It is expected to cover:

- Leave type definitions
- Employee leave balance visibility
- Leave request submission
- Manager approval or rejection
- Balance deduction after approved leave

From a backend perspective, this module owns leave requests, approval rules, balance updates, and related validations.

## Product and Architecture Direction

The project starts as a `modular monolith`.

That means:

- One deployable backend application
- One primary relational database
- Clear module boundaries in code
- Independent business logic per module without the operational cost of microservices

This is a deliberate choice for v1. It keeps implementation and deployment practical while still encouraging separation of concerns, maintainability, and future extraction if the product outgrows a single service.

## Technology Stack

The initial stack is intentionally conservative and product-oriented:

- `Go`
  Chosen for backend performance, clear concurrency primitives, a strong standard library, and operational simplicity.
- `PostgreSQL`
  Chosen for relational data integrity, rich querying, and suitability for HR workflows that involve linked entities and reporting.
- `GORM`
  Chosen to accelerate development of the initial data access layer while still allowing direct SQL when domain-specific queries become more complex.
- `REST JSON`
  Chosen as the first external interface because it is easy to integrate with both future web and mobile clients.
- `JWT`
  Chosen for straightforward stateless authentication across multiple clients.

## Why These Decisions Fit v1

- `Go` keeps the service efficient and operationally simple.
- `PostgreSQL` fits structured HR data better than document-first storage.
- `GORM` helps move quickly without committing the system to fully generic abstractions.
- `REST` avoids premature API complexity.
- `JWT` works well for both browser-based and mobile clients.
- A modular monolith keeps development fast while preserving domain boundaries.

## High-Level Architecture

At a high level, the system is expected to follow this request path:

```text
Client (future web admin / future Flutter app / API consumer)
    ->
HTTP Router + Middleware
    ->
Auth / Authorization checks
    ->
Module Handler
    ->
Module Service / Business Logic
    ->
Repository / Data Access
    ->
PostgreSQL
```

Supporting cross-cutting concerns will live in shared platform packages:

- configuration loading
- database initialization
- authentication helpers
- middleware
- logging
- health checks

## Proposed Repository Structure

The repository is currently empty, but the intended backend layout is:

```text
.
├── cmd/
│   └── api/
├── internal/
│   ├── platform/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── db/
│   │   ├── http/
│   │   └── logging/
│   └── modules/
│       ├── corehr/
│       └── leave/
├── migrations/
├── docs/
│   └── openapi/
├── scripts/
└── README.md
```

### Directory Responsibilities

- `cmd/api`
  Application entrypoint, startup wiring, and dependency composition.
- `internal/platform`
  Shared infrastructure code that is not domain-specific.
- `internal/modules/corehr`
  Employee records, organization relationships, and asset management logic.
- `internal/modules/leave`
  Leave requests, balances, approval flow, and leave rules.
- `migrations`
  Database schema evolution owned by the backend codebase.
- `docs/openapi`
  Generated or maintained API documentation artifacts.
- `scripts`
  Local developer helpers such as setup and workflow scripts.

## Domain Model Overview

The exact schema is intentionally not frozen yet, but the initial domain language is already clear.

### Identity and Access

- `users`
  Application login accounts and role membership.

### Core HR

- `employees`
  Employee profile and employment-related information.
- `departments`
  Department definitions used for organization grouping.
- `positions`
  Position or title definitions attached to employees.
- `assets`
  Company-owned equipment such as laptops and phones.
- `asset_assignments`
  Assignment history and current asset responsibility.

### Leave & Attendance

- `leave_types`
  Definitions such as annual leave or sick leave.
- `leave_balances`
  Employee leave totals, used amounts, and remaining balances.
- `leave_requests`
  Submitted leave requests and their approval state.

### Organizational Relationship

The hierarchy is expected to be modeled through the employee-to-manager relationship, allowing the system to answer:

- who a person reports to
- who reports to a manager
- which manager should review a leave request

## Authentication and Authorization

The initial authentication model is:

- email/password sign-in
- JWT-based authorization for API access

The initial authorization model is role-based.

### Roles

- `Admin`
  Full system-level access and operational control.
- `HR`
  Employee record management, asset oversight, and HR-oriented leave visibility.
- `Manager`
  Team-level visibility and leave approval responsibility for direct reports.
- `Employee`
  Self-service access to personal and leave-related data within policy limits.

### Role Permission Matrix

| Capability | Admin | HR | Manager | Employee |
| --- | --- | --- | --- | --- |
| View all employees | Yes | Yes | Team-focused | No |
| Edit employee records | Yes | Yes | No | No |
| View own profile | Yes | Yes | Yes | Yes |
| View team structure | Yes | Yes | Yes | Limited |
| Manage assets | Yes | Yes | No | No |
| Submit leave request | Yes | Yes | Yes | Yes |
| Approve leave requests | Yes | HR policy-dependent | Yes | No |
| View own leave balance | Yes | Yes | Yes | Yes |

This matrix is directional guidance for v1 and should be refined in code once implementation begins.

## API Design Principles

The backend will expose a versioned REST interface under:

```text
/api/v1
```

### Resource Groups

The main API groups planned for v1 are:

- `auth`
- `employees`
- `organization`
- `assets`
- `leave-types`
- `leave-balances`
- `leave-requests`

### API Conventions

- JSON request and response bodies
- Resource-oriented route naming
- Clear separation between read and mutation endpoints
- Explicit authentication and authorization checks at the API boundary
- A documented error response structure, to be finalized during implementation
- OpenAPI documentation planned as part of the backend foundation

### Example Endpoint Groups

The following examples illustrate the intended API surface without locking final payload shapes too early:

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

## Key Backend Workflows

### Employee Creation and Update

1. An authorized user creates or updates an employee record.
2. Core identity and employment metadata are validated.
3. The employee is linked to a department, position, and optionally a manager.
4. The record becomes available to dependent modules such as leave approvals and organization views.

### Organization Hierarchy Lookup

1. The system reads employee-to-manager relationships.
2. The API returns a tree or structured list representation.
3. Future clients can render this as an org chart without embedding hierarchy rules in the UI.

### Asset Assignment and Return

1. HR or Admin registers an asset.
2. The asset is assigned to an employee.
3. The system tracks status, assignment date, and return state.
4. Returned assets remain historically traceable instead of disappearing from records.

### Leave Request Submission

1. An employee submits a leave request.
2. The system validates dates, type, and available balance.
3. The request is routed to the employee's direct manager.
4. The request enters a pending approval state.

### Manager Approval or Rejection

1. The direct manager reviews pending requests.
2. The manager approves or rejects the request.
3. Approval updates the request status and triggers balance deduction.
4. Rejection stores a reason when provided and preserves balance.

### Balance Deduction After Approval

1. The system calculates the approved leave duration based on the policy in use.
2. Remaining balance is reduced only after approval.
3. Balance state remains queryable for both the employee and authorized HR roles.

## Local Development Setup

Implementation has not started yet, but the intended local development flow is straightforward.

### Current Repository Status

- The repository is currently a foundation workspace.
- No Go module or application code has been created yet.
- No migration tooling has been committed yet.
- No database bootstrap scripts exist yet.

### Prerequisites

Before implementation begins, a local developer machine should have:

- `Go`
- `PostgreSQL`
- `Git`
- a shell environment such as `zsh` or `bash`

### Important Note About This Machine

At the time this README was prepared, the `go` command was not available in the current environment. That means Go must be installed before backend implementation can start successfully on this machine.

### Intended First-Time Setup Flow

Once implementation begins, the expected onboarding flow should look roughly like this:

1. Install Go.
2. Install or start PostgreSQL locally.
3. Clone the repository.
4. Create a local environment file from an example template.
5. Create the development database.
6. Run database migrations.
7. Start the API server.

Exact commands will be added after the codebase and toolchain decisions are committed to the repository.

## Configuration

The backend should be configured through environment variables with a single typed configuration layer in code.

### Expected Configuration Categories

- application environment
- HTTP port
- database host, port, name, user, password, SSL mode
- JWT secret or signing configuration
- token lifetime settings
- logging level

### Example Variables

These names are illustrative and may be refined during implementation:

```text
APP_ENV=development
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=hr_project
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
JWT_SECRET=change-me
JWT_ACCESS_TOKEN_TTL=15m
```

Configuration should be loaded centrally and passed into startup wiring rather than read ad hoc across the codebase.

## Database and Migrations

`PostgreSQL` is the planned primary database.

### Migration Ownership

- Schema migrations should live in the repository.
- Backend changes that require schema updates should include corresponding migrations.
- Migration history should be deterministic and safe to apply in order.

### Seed Data Expectations

Initial seed data will likely be useful for:

- base roles
- leave types
- a small demo org structure
- a sample admin user for local development

Seed strategy should stay simple in v1 and avoid hiding business rules in ad hoc scripts.

### Local Database Bootstrap Flow

The intended local flow is:

1. create database
2. apply migrations
3. load minimal seed data
4. start API

The exact migration tool is intentionally not locked in this README yet. That decision should be made once implementation begins so the team can pick tooling that matches the project style and operational needs.

## Testing Strategy

The backend should be designed with multiple levels of validation.

### Unit Tests

Best for:

- leave balance calculations
- approval rule checks
- authorization helpers
- employee-manager relationship logic

### Integration Tests

Best for:

- authentication flow
- database-backed employee creation
- asset assignment and return
- leave request creation
- manager approval and balance update behavior

### API Contract Tests

Best for:

- preserving route behavior
- validating request and response shape expectations
- reducing breaking changes for future clients

### Example Critical Test Scenarios

- an employee cannot approve their own leave request unless policy explicitly allows it
- a manager cannot approve requests for employees outside their reporting line
- a leave request should not reduce balance before approval
- a returned asset should no longer appear as an active assignment
- an employee should not be able to edit another employee's private record

## Engineering Conventions

The project should favor clarity over cleverness.

### Package Boundaries

- shared infrastructure belongs in `internal/platform`
- business rules belong inside their owning module
- domain logic should not leak into generic middleware or utility packages
- repositories should remain focused on persistence, not policy decisions

### Branching

- use feature branches
- default branch naming can follow the `codex/` prefix convention when created through the local tooling context
- branch names should describe the unit of work clearly

### Commits

- keep commits focused
- avoid mixing refactors with behavior changes when possible
- prefer small, reviewable increments

### Repository Hygiene

- do not commit secrets
- keep generated files intentional
- document architectural decisions that materially affect module boundaries
- prefer explicit naming over ambiguous helpers

## Roadmap and Future Phases

The backend foundation is intended to support several next steps after v1 starts taking shape:

- Flutter mobile app for employee self-service
- Web admin interface for HR and managers
- Notification workflows for approvals and reminders
- More advanced leave policy rules
- Multi-tenant readiness for SaaS expansion
- Payroll-adjacent integrations
- Reporting and dashboards
- Audit trails and compliance-oriented event history

## What This README Should Help a New Developer Understand

After reading this document, a new developer should be able to answer:

- What problem does this project solve?
- What is included in v1?
- What is not built yet?
- Why is the backend being built with Go and PostgreSQL?
- What are the main modules?
- How will authentication and authorization work?
- How does the leave approval flow work?
- What will the repository structure look like?
- What will local setup require once implementation begins?

## Status

This repository is currently in the planning and foundation stage. The README reflects agreed architectural and product decisions, but it intentionally avoids pretending that unresolved implementation details are already final. As the codebase is created, this document should evolve alongside the actual project structure and developer workflow.
