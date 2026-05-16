# Implementation Plan

## Scope

Split `projex_insights_handlers.go` (904 lines) and `projex_insights_handlers_test.go` (901 lines) into 6+6 domain files.

## Task DAG

1. Split handlers file — source: proposal
   - Dependencies: none
   - Validation: `go build ./pkg/toolset/yunxiao`
   - Create 6 new handler files, delete original

2. Split test file — source: proposal
   - Dependencies: task 1
   - Validation: `go test ./pkg/toolset/yunxiao`
   - Create 6 new test files, delete original

3. Verify — source: proposal
   - Dependencies: tasks 1-2
   - Validation: `make lint`, `make test`, `make coverage-check`

## Handler Split Mapping

From `projex_insights_handlers.go`:

**File 1: projex_insights_risk_handlers.go**
- handleGetProjectRiskDashboard (lines 15-44)
- addRiskCategorySections (lines 83-92)
- addRiskFocusSections (lines 94-122)
- riskDashboardFilters (lines 244-260)

**File 2: projex_insights_member_handlers.go**
- handleGetProjectMemberTaskStatus (lines 46-81)
- projectMemberTaskPayload (lines 124-153)
- projectMemberStatusGroups (lines 155-172)
- projectTaskStatusMembers (lines 174-188)
- projectMembersFromResponse (lines 190-213)
- searchProjectWorkitems (lines 215-223)
- parseStatusGroups (lines 225-242)
- projectTaskStatusFilters (lines 262-279)
- copyParams (lines 281-287)
- todayDate (lines 289-291)

**File 3: projex_insights_velocity_handlers.go**
- handleGetSprintVelocity (lines 295-343)
- clampSprintCount (lines 345-353)
- parseSprintList (lines 355-365)
- buildSprintVelocityStats (lines 367-406)
- countCompletedWorkitems (lines 408-420)
- sprintVelocityFilters (lines 422-428)

**File 4: projex_insights_timeline_handlers.go**
- handleGetWorkitemStatusTimeline (lines 430-473)
- workitemStatusTimelineFilters (lines 475-480)
- parseStatusTimeline (lines 482-519)
- calculateTimelineStats (lines 521-539)

**File 5: projex_insights_blocker_handlers.go**
- handleGetBlockerAnalysis (lines 541-603)
- fetchCategoryWorkitems (lines 605-619)
- checkWorkitemBlockers (lines 621-651)
- parseAnyList (lines 653-663)
- blockerAnalysisFilters (lines 665-670)

**File 6: projex_insights_workload_handlers.go**
- handleGetMemberWorkloadTrend (lines 672-728)
- buildMemberWorkload (lines 730-763)
- countWorkitemStatuses (lines 765-785)
- memberWorkloadTrendFilters (lines 787-793)
- handleGetTeamWorkloadBreakdown (lines 795-873)
- teamWorkloadBreakdownFilters (lines 875-882)
- extractWorkitemStatusName (lines 884-891)
- extractLabelNames (lines 893-904)

## Test Split Mapping

Each test file mirrors the handler split — tests for functions in file N go to the corresponding `_test.go` file.

## Imports

All split files remain in package `yunxiao`. Each file gets only the imports it needs (minimal — `goimports` or manual check). Most files need only `context`, `fmt`, `net/http`, `net/url`, and `strconv` subsets.

## Review Status

self-reviewed: zero P0/P1. Pure mechanical split, no logic changes.
