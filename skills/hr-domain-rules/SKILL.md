---
name: hr-domain-rules
description: Use this skill when implementing HR business rules for this repository, especially employee records, reporting lines, asset assignment, leave requests, approvals, and balance updates. It defines the agreed v1 product rules and constraints.
---

# HR Domain Rules

Use this skill when translating product requirements into backend rules for the HR system.

## Product Scope

This project currently focuses on two modules:

- `Core HR`
- `Leave & Attendance`

Assume `single-company v1` unless the task explicitly introduces multi-tenant behavior.

## Core HR Rules

- Every employee record represents one worker in the company.
- An employee may optionally have a direct manager.
- Manager relationships drive approval visibility in the leave module.
- Asset records must support assignment history, not just current ownership.
- Returned assets should remain queryable historically.

## Leave Rules

- Leave requests are created by the employee.
- Approval is single-step in v1.
- The approver is the employee's direct manager by default.
- Leave balance should only be reduced after approval.
- Rejected requests must not consume leave balance.
- If balance is insufficient, the request should be blocked or marked invalid according to the use case being implemented.

## Authorization Defaults

- `Employee` can access self-service data.
- `Manager` can review direct-report leave actions.
- `HR` can manage records needed to operate the HR system.
- `Admin` has full system access.

## Modeling Guidance

- Prefer explicit fields for sensitive HR data instead of collapsing everything into generic JSON.
- Keep status fields narrow and meaningful.
- Preserve audit-friendly history where the business process matters, especially for assets and leave approvals.

## Change Safety

When a task proposes a new workflow, check whether it changes any of these agreed defaults:

- single-company scope
- single-level manager approval
- role model of `Admin`, `HR`, `Manager`, `Employee`
- API-first backend direction

If yes, call that out clearly before implementing.

## When To Read References

- Read [Domain Glossary](references/domain-glossary.md) when naming entities, statuses, or flows.
