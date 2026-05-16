package yunxiao

import (
	"context"
	"net/url"
)

func handleGetRepositoryOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	repoPath := codeupRepositoryPath(organizationID, repositoryID)

	repository, err := getProjectOverviewSection(ctx, c, "repository", repoPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"repository": repository,
		"filters":    repositoryOverviewFilters(params),
	}

	branches, err := fetchBranches(ctx, c, params, repoPath)
	if err != nil {
		return "", err
	}
	if branches != nil {
		overview["branches"] = branches
	}

	refName := resolveRefName(params, repository)

	commits, err := fetchCommits(ctx, c, params, repoPath, refName)
	if err != nil {
		return "", err
	}
	if commits != nil {
		overview["commits"] = commits
	}

	mrs, err := fetchMergeRequests(ctx, c, params, organizationID, repositoryID)
	if err != nil {
		return "", err
	}
	if mrs != nil {
		overview["mergeRequests"] = mrs
	}

	return marshalPretty(overview)
}

func fetchBranches(ctx context.Context, c *Client, params map[string]any, repoPath string) (any, error) {
	if !optionalBoolDefault(params, "includeBranches", true) {
		return nil, nil
	}
	return getProjectOverviewSection(ctx, c, "branches", repoPath+"/branches", pageOneLimitQuery(params, "branchLimit", 5))
}

func resolveRefName(params map[string]any, repository any) string {
	if refName := optionalStringDefault(params, "refName", ""); refName != "" {
		return refName
	}
	if repoMap, ok := repository.(map[string]any); ok {
		if db, ok := repoMap["defaultBranch"].(string); ok {
			return db
		}
	}
	return ""
}

func fetchCommits(ctx context.Context, c *Client, params map[string]any, repoPath, refName string) (any, error) {
	if !optionalBoolDefault(params, "includeCommits", true) || refName == "" {
		return nil, nil
	}
	commitQuery := pageOneLimitQuery(params, "commitLimit", 5)
	commitQuery.Set("refName", refName)
	return getProjectOverviewSection(ctx, c, "commits", repoPath+"/commits", commitQuery)
}

func fetchMergeRequests(ctx context.Context, c *Client, params map[string]any, organizationID, repositoryID string) (any, error) {
	if !optionalBoolDefault(params, "includeMergeRequests", true) {
		return nil, nil
	}
	mrQuery := pageOneLimitQuery(params, "mrLimit", 5)
	mrQuery.Set("state", optionalStringDefault(params, "mrState", "opened"))
	mrQuery.Add("repositoryIds", repositoryID)
	return getProjectOverviewSection(ctx, c, "mergeRequests", codeupOrganizationPath(organizationID)+"/mergeRequests", mrQuery)
}

func repositoryOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeBranches":      optionalBoolDefault(params, "includeBranches", true),
		"includeCommits":       optionalBoolDefault(params, "includeCommits", true),
		"includeMergeRequests": optionalBoolDefault(params, "includeMergeRequests", true),
		"refName":              optionalStringDefault(params, "refName", ""),
		"branchLimit":          optionalIntDefault(params, "branchLimit", 5),
		"commitLimit":          optionalIntDefault(params, "commitLimit", 5),
		"mrLimit":              optionalIntDefault(params, "mrLimit", 5),
		"mrState":              optionalStringDefault(params, "mrState", "opened"),
	}
}

func handleGetChangeRequestOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	crPath := changeRequestPath(organizationID, repositoryID, localID)

	cr, err := c.GetJSON(ctx, crPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"changeRequest": cr,
		"filters":       changeRequestOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includePatchSets", true) {
		patches, err := c.GetJSON(ctx, crPath+"/diffs/patches", nil)
		if err != nil {
			return "", err
		}
		overview["patchSets"] = patches
	}

	if optionalBoolDefault(params, "includeComments", true) {
		body := map[string]any{
			"comment_type": "GLOBAL_COMMENT",
			"state":        optionalStringDefault(params, "commentState", "OPENED"),
			"resolved":     optionalBoolDefault(params, "commentResolved", false),
		}
		comments, err := c.PostJSONWithMetadata(ctx, crPath+"/comments/list", body)
		if err != nil {
			return "", err
		}
		overview["comments"] = comments
	}

	return marshalPretty(overview)
}

func changeRequestOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includePatchSets": optionalBoolDefault(params, "includePatchSets", true),
		"includeComments":  optionalBoolDefault(params, "includeComments", true),
		"commentState":     optionalStringDefault(params, "commentState", "OPENED"),
		"commentResolved":  optionalBoolDefault(params, "commentResolved", false),
	}
}

func handleGetCommitOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	sha, err := requiredString(params, "sha")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	repoPath := codeupRepositoryPath(organizationID, repositoryID)
	commitPath := repoPath + "/commits/" + url.PathEscape(sha)

	commit, err := getProjectOverviewSection(ctx, c, "commit", commitPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"commit":  commit,
		"filters": commitOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeStatuses", true) {
		statusQuery := pageOneLimitQuery(params, "statusLimit", 5)
		statuses, err := getProjectOverviewSection(ctx, c, "statuses", commitPath+"/statuses", statusQuery)
		if err != nil {
			return "", err
		}
		overview["statuses"] = statuses
	}

	if optionalBoolDefault(params, "includeCheckRuns", true) {
		checkRunQuery := pageOneLimitQuery(params, "checkRunLimit", 5)
		checkRunQuery.Set("ref", sha)
		checkRuns, err := getProjectOverviewSection(ctx, c, "checkRuns", repoPath+"/checkRuns", checkRunQuery)
		if err != nil {
			return "", err
		}
		overview["checkRuns"] = checkRuns
	}

	return marshalPretty(overview)
}

func commitOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeStatuses":  optionalBoolDefault(params, "includeStatuses", true),
		"includeCheckRuns": optionalBoolDefault(params, "includeCheckRuns", true),
		"statusLimit":      optionalIntDefault(params, "statusLimit", 5),
		"checkRunLimit":    optionalIntDefault(params, "checkRunLimit", 5),
	}
}

func handleGetBranchOverview(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	branchName, err := requiredString(params, "branchName")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	repoPath := codeupRepositoryPath(organizationID, repositoryID)
	branchPath := repoPath + "/branches/" + url.PathEscape(branchName)

	branch, err := getProjectOverviewSection(ctx, c, "branch", branchPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"branch":  branch,
		"filters": branchOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeCommits", true) {
		commitQuery := pageOneLimitQuery(params, "commitLimit", 5)
		commitQuery.Set("refName", branchName)
		commits, err := getProjectOverviewSection(ctx, c, "commits", repoPath+"/commits", commitQuery)
		if err != nil {
			return "", err
		}
		overview["commits"] = commits
	}

	if optionalBoolDefault(params, "includeMergeRequests", true) {
		mrQuery := pageOneLimitQuery(params, "mrLimit", 5)
		mrQuery.Set("state", optionalStringDefault(params, "mrState", "opened"))
		mrQuery.Add("repositoryIds", repositoryID)
		mrQuery.Set("targetBranch", branchName)
		mrs, err := getProjectOverviewSection(ctx, c, "mergeRequests", codeupOrganizationPath(organizationID)+"/mergeRequests", mrQuery)
		if err != nil {
			return "", err
		}
		overview["mergeRequests"] = mrs
	}

	return marshalPretty(overview)
}

func branchOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeCommits":       optionalBoolDefault(params, "includeCommits", true),
		"includeMergeRequests": optionalBoolDefault(params, "includeMergeRequests", true),
		"commitLimit":          optionalIntDefault(params, "commitLimit", 5),
		"mrLimit":              optionalIntDefault(params, "mrLimit", 5),
		"mrState":              optionalStringDefault(params, "mrState", "opened"),
	}
}
