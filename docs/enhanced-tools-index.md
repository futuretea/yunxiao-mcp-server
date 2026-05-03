# Enhanced Tools Index

This project provides **enhanced aggregation tools** that combine multiple Yunxiao OpenAPI calls into single, user-centric operations. When both a base tool and an enhanced tool exist for the same concept, prefer the enhanced tool — it returns aggregated, filtered responses that are easier to summarize.

## Projex (Project Management)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_project_overview` | `get_project` + members + sprints + milestones + versions + labels | Quick project snapshot |
| `get_project_workitem_summary` | `search_workitems` per category | Work item shape by category |
| `get_project_workitem_context` | work item types + members + labels + optional fields/workflow | Preparing to create/edit a work item |
| `get_sprint_overview` | `get_sprint` + `search_workitems` per category | Sprint progress tracking |
| `get_my_project_workitems` | `search_workitems` per category filtered by user | Personal task dashboard |
| `get_project_workitem_board` | `search_workitems` grouped by status | Kanban-style board view |
| `get_project_workitem_detail` | `get_workitem` + activities + attachments + comments + relations | Complete single work item view |
| `get_project_risk_dashboard` | `search_workitems` for overdue + high-priority + stale | Risk identification |
| `get_project_member_task_status` | `search_workitems` per member + status groups | Workload distribution |

See [`projex-enhanced-tools.md`](projex-enhanced-tools.md) for detailed parameters and examples.

## Codeup (Source Code Hosting)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_repository_overview` | `get_repository` + branches + commits + merge requests | Repository snapshot |

See [`codeup-enhanced-tools.md`](codeup-enhanced-tools.md) for detailed parameters and examples.

## Flow (CI/CD Pipeline)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_pipeline_overview` | `get_pipeline` + latest run + recent run history | Pipeline status check |

See [`flow-enhanced-tools.md`](flow-enhanced-tools.md) for detailed parameters and examples.

## Appstack (Application Deployment)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_application_overview` | `get_application` + environments + orchestrations | Application snapshot |

See [`appstack-enhanced-tools.md`](appstack-enhanced-tools.md) for detailed parameters and examples.

## Platform (Organization Management)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_organization_overview` | `get_organization` + departments + members + groups + roles | Organization snapshot |

See [`platform-enhanced-tools.md`](platform-enhanced-tools.md) for detailed parameters and examples.

## Quick Reference

**Total enhanced tools**: 13 across 5 domains.

**Base tools they replace**: ~31 raw API calls.

**Common patterns**:
- Most enhanced tools accept `organizationId` (optional, auto-injected) and a domain-specific ID.
- Section toggles (`includeXxx`) let you control API load.
- Limit parameters (`sampleLimit`, `branchLimit`, `runLimit`, `envLimit`) default to small values (5) to keep responses concise.
