# Conditions Cookbook

Yunxiao uses a JSON `conditions` object for advanced filtering in search endpoints. Several MCP tools expose a `conditions` parameter that accepts a raw JSON string. When the parameter is omitted, the server automatically builds conditions from simple filter parameters (`subject`, `status`, `assignedTo`, etc.).

This document provides common `conditions` JSON examples for tools that support them.

## When to Use Raw Conditions

- You need a filter combination that is **not** exposed as a simple parameter.
- You need an operator other than `CONTAINS` or `BETWEEN` (e.g., `EQUALS`, `NOT_CONTAINS`).
- You are debugging a query and want full control over the filter shape.

For everyday queries, prefer the simple filter parameters — the server translates them into conditions automatically.

## Conditions JSON Format

A conditions object has one top-level key:

```json
{
  "conditionGroups": [
    [
      { "fieldIdentifier": "...", "operator": "...", "className": "...", "format": "...", "value": [...] }
    ]
  ]
}
```

- `conditionGroups` is an array of filter groups. Each group is an array of condition objects.
- Conditions inside the same group are combined with **AND** logic.
- `fieldIdentifier`: the field to filter on.
- `operator`: `CONTAINS`, `BETWEEN`, `EQUALS`, etc.
- `className`: `string`, `status`, `user`, `sprint`, `workitemType`, `option`, `dateTime`, `date`, `tag`.
- `format`: `input`, `list`, `multiList`.
- `value`: array of strings.
- `toValue`: used with `BETWEEN` for the upper bound; otherwise `null`.

## search_workitems

### Subject contains a keyword

```json
{
  "conditionGroups": [
    [
      {
        "className": "string",
        "fieldIdentifier": "subject",
        "format": "input",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["login"]
      }
    ]
  ]
}
```

Equivalent simple filter: `subject=login`

### Status is TODO or DOING

```json
{
  "conditionGroups": [
    [
      {
        "className": "status",
        "fieldIdentifier": "status",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["TODO", "DOING"]
      }
    ]
  ]
}
```

Equivalent simple filter: `status=TODO,DOING`

### Assigned to specific users

```json
{
  "conditionGroups": [
    [
      {
        "className": "user",
        "fieldIdentifier": "assignedTo",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["user-123", "user-456"]
      }
    ]
  ]
}
```

Equivalent simple filter: `assignedTo=user-123,user-456`

### Created by a specific user

```json
{
  "conditionGroups": [
    [
      {
        "className": "user",
        "fieldIdentifier": "creator",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["user-123"]
      }
    ]
  ]
}
```

Equivalent simple filter: `creator=user-123`

### In a specific sprint

```json
{
  "conditionGroups": [
    [
      {
        "className": "sprint",
        "fieldIdentifier": "sprint",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["sprint-789"]
      }
    ]
  ]
}
```

Equivalent simple filter: `sprint=sprint-789`

### Work item type filter

```json
{
  "conditionGroups": [
    [
      {
        "className": "workitemType",
        "fieldIdentifier": "workitemType",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["type-abc"]
      }
    ]
  ]
}
```

Equivalent simple filter: `workitemType=type-abc`

### Priority filter

```json
{
  "conditionGroups": [
    [
      {
        "className": "option",
        "fieldIdentifier": "priority",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["high", "urgent"]
      }
    ]
  ]
}
```

Equivalent simple filter: `priority=high,urgent`

### Tag filter (multiList)

```json
{
  "conditionGroups": [
    [
      {
        "className": "tag",
        "fieldIdentifier": "tag",
        "format": "multiList",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["tag-1", "tag-2"]
      }
    ]
  ]
}
```

Equivalent simple filter: `tag=tag-1,tag-2`

### Created within a date range

```json
{
  "conditionGroups": [
    [
      {
        "className": "dateTime",
        "fieldIdentifier": "gmtCreate",
        "format": "input",
        "operator": "BETWEEN",
        "toValue": "2026-04-30 23:59:59",
        "value": ["2026-04-01 00:00:00"]
      }
    ]
  ]
}
```

Equivalent simple filters: `createdAfter=2026-04-01&createdBefore=2026-04-30`

### Combined: subject + status + assignedTo

```json
{
  "conditionGroups": [
    [
      {
        "className": "string",
        "fieldIdentifier": "subject",
        "format": "input",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["login"]
      },
      {
        "className": "status",
        "fieldIdentifier": "status",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["TODO", "DOING"]
      },
      {
        "className": "user",
        "fieldIdentifier": "assignedTo",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["user-123"]
      }
    ]
  ]
}
```

This combines all three filters with AND logic.

## search_projects

### Project name contains a keyword

```json
{
  "conditionGroups": [
    [
      {
        "className": "string",
        "fieldIdentifier": "name",
        "format": "input",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["backend"]
      }
    ]
  ]
}
```

Equivalent simple filter: `name=backend`

### Filter by project status

```json
{
  "conditionGroups": [
    [
      {
        "className": "status",
        "fieldIdentifier": "status",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["active"]
      }
    ]
  ]
}
```

Equivalent simple filter: `status=active`

### Filter by creator

```json
{
  "conditionGroups": [
    [
      {
        "className": "user",
        "fieldIdentifier": "creator",
        "format": "list",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["user-123"]
      }
    ]
  ]
}
```

Equivalent simple filter: `creator=user-123`

## search_testcases

### Subject contains a keyword

```json
{
  "conditionGroups": [
    [
      {
        "className": "string",
        "fieldIdentifier": "subject",
        "format": "input",
        "operator": "CONTAINS",
        "toValue": null,
        "value": ["login"]
      }
    ]
  ]
}
```

Equivalent simple filter: `subject=login`

## Simple Filter to Conditions Mapping

| Simple Parameter | Field Identifier | Class Name | Format |
|------------------|------------------|------------|--------|
| `subject` | `subject` | `string` | `input` |
| `status` | `status` | `status` | `list` |
| `assignedTo` | `assignedTo` | `user` | `list` |
| `creator` | `creator` | `user` | `list` |
| `sprint` | `sprint` | `sprint` | `list` |
| `workitemType` | `workitemType` | `workitemType` | `list` |
| `statusStage` | `statusStage` | `statusStage` | `list` |
| `priority` | `priority` | `option` | `list` |
| `tag` | `tag` | `tag` | `multiList` |
| `subjectDescription` | `subject-description` | `string` | `input` |
| `createdAfter`/`createdBefore` | `gmtCreate` | `dateTime` | `input` |
| `updatedAfter`/`updatedBefore` | `gmtModified` | `dateTime` | `input` |
| `finishTimeAfter`/`finishTimeBefore` | `finishTime` | `date` | `input` |
| `updateStatusAtAfter`/`updateStatusAtBefore` | `updateStatusAt` | `date` | `input` |

## Tips

- Pass `conditions` as a **single-line JSON string** (no newlines) in the MCP tool parameter.
- When `conditions` is provided, simple filter parameters are **ignored** by the server for that request.
- Date range `BETWEEN` operators use `YYYY-MM-DD HH:MM:SS` format in the `value` and `toValue` fields.
- Use the enhanced tools (`get_project_workitem_summary`, `get_sprint_overview`) when you do not need raw conditions — they translate simple parameters automatically.
