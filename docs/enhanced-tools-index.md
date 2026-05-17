# Enhanced Tools Index

This project provides **enhanced aggregation tools** that combine multiple Yunxiao OpenAPI calls into single, user-centric operations. When both a base tool and an enhanced tool exist for the same concept, prefer the enhanced tool — it returns aggregated, filtered responses that are easier to summarize.

## Projex (Project Management)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_project_overview` | project info + members + sprints + milestones + versions + labels | Quick project snapshot |
| `get_project_workitem_summary` | `search_workitems` per category | Work item shape by category |
| `get_project_workitem_context` | work item types + members + labels + optional fields/workflow | Preparing to create/edit a work item |
| `get_sprint_overview` | sprint info + `search_workitems` per category | Sprint progress tracking |
| `get_my_project_workitems` | `search_workitems` per category filtered by user | Personal task dashboard |
| `get_project_workitem_board` | `search_workitems` grouped by status | Kanban-style board view |
| `get_project_workitem_detail` | work item info + activities + attachments + comments + relations | Complete single work item view |
| `get_project_risk_dashboard` | `search_workitems` for overdue + high-priority + stale | Risk identification |
| `get_project_member_task_status` | `search_workitems` per member + status groups | Workload distribution |
| `get_work_item_type_overview` | work item type info + field config + workflow | Work item type configuration |

See [`projex-enhanced-tools.md`](projex-enhanced-tools.md) for detailed parameters and examples.

## Codeup (Source Code Hosting)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_repository_overview` | repository info + branches + commits + merge requests | Repository snapshot |
| `get_change_request_overview` | change request info + patch sets + comments | Change request snapshot |
| `get_commit_overview` | commit info + commit statuses + check runs | Commit CI/CD health snapshot |
| `get_branch_overview` | branch info + recent commits + merge requests | Branch activity snapshot |

See [`codeup-enhanced-tools.md`](codeup-enhanced-tools.md) for detailed parameters and examples.

## Flow (CI/CD Pipeline)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_pipeline_overview` | pipeline info + latest run + recent run history | Pipeline status check |
| `get_pipeline_run_overview` | pipeline run info + jobs by category | Pipeline run details |

See [`flow-enhanced-tools.md`](flow-enhanced-tools.md) for detailed parameters and examples.

## Appstack (Application Deployment)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_application_overview` | application info + environments + orchestrations | Application snapshot |
| `get_environment_overview` | environment info + variable groups + latest orchestration | Environment snapshot |
| `get_release_overview` | release info + members + products + change requests | Release snapshot |
| `get_system_overview` | system info + attached apps + members | System snapshot |
| `get_app_release_stage_overview` | stage info + pipeline run + integrated metadata | Stage execution debugging |
| `get_change_order_overview` | change order info + jobs | Deployment order details |

See [`appstack-enhanced-tools.md`](appstack-enhanced-tools.md) for detailed parameters and examples.

## Platform (Organization Management)

| Tool | What it combines | When to use |
|------|-----------------|-------------|
| `get_organization_overview` | organization info + departments + members + groups + roles | Organization snapshot |
| `get_organization_department_overview` | department info + ancestors | Department snapshot |
| `get_organization_group_overview` | group info + members | Group snapshot |

See [`platform-enhanced-tools.md`](platform-enhanced-tools.md) for detailed parameters and examples.

## Quick Reference

**Total enhanced tools**: 25 across 5 domains.

Use `--compact` to hide the 28 raw tools superseded by these enhanced tools.

**Common patterns**:
- Most enhanced tools accept `organizationId` (optional, auto-injected) and a domain-specific ID.
- Section toggles (`includeXxx`) let you control API load.
- Limit parameters (`sampleLimit`, `branchLimit`, `runLimit`, `envLimit`) default to small values (5) to keep responses concise.
