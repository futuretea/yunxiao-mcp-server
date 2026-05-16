# Split projex_insights_handlers.go

## Summary

`pkg/toolset/yunxiao/projex_insights_handlers.go` is 904 lines (37 functions), exceeding the 600-line hard limit by 50%. Split it into 6 domain files by logical group, and split the corresponding 901-line test file to match.

## Motivation

### Goals

- Each new file < 250 LOC
- No behavior changes — pure mechanical split
- `make lint` + `make test` pass, no coverage regression
- Test file split mirrors handler file split

### Non-Goals

- No logic refactoring, renaming, or reordering
- No new tests or test changes beyond file placement
- No import optimization beyond what each split file needs

## Proposal

### Split Plan

| New File | Functions | Est. LOC |
|----------|-----------|----------|
| `projex_insights_risk_handlers.go` | handleGetProjectRiskDashboard, addRiskCategorySections, addRiskFocusSections, riskDashboardFilters | ~100 |
| `projex_insights_member_handlers.go` | handleGetProjectMemberTaskStatus, projectMemberTaskPayload, projectMemberStatusGroups, projectTaskStatusMembers, projectMembersFromResponse, searchProjectWorkitems, parseStatusGroups, projectTaskStatusFilters, copyParams, todayDate | ~175 |
| `projex_insights_velocity_handlers.go` | handleGetSprintVelocity, clampSprintCount, parseSprintList, buildSprintVelocityStats, countCompletedWorkitems, sprintVelocityFilters | ~150 |
| `projex_insights_timeline_handlers.go` | handleGetWorkitemStatusTimeline, workitemStatusTimelineFilters, parseStatusTimeline, calculateTimelineStats | ~110 |
| `projex_insights_blocker_handlers.go` | handleGetBlockerAnalysis, fetchCategoryWorkitems, checkWorkitemBlockers, parseAnyList, blockerAnalysisFilters | ~125 |
| `projex_insights_workload_handlers.go` | handleGetMemberWorkloadTrend, buildMemberWorkload, countWorkitemStatuses, memberWorkloadTrendFilters, handleGetTeamWorkloadBreakdown, teamWorkloadBreakdownFilters, extractWorkitemStatusName, extractLabelNames | ~220 |

### Why 6 files not 5

Original plan grouped blocker + workload together, but combined they exceed 250 LOC. Splitting into 6 keeps each file under the target.

### Risk Assessment

- **Risk**: Low. Pure code movement within the same package. No import changes needed (all files share the same package `yunxiao`).
- **Mitigation**: `make test` before and after verifies zero behavior change.

## Review Status

self-reviewed: zero P0/P1. Mechanical refactoring, no design decisions.
