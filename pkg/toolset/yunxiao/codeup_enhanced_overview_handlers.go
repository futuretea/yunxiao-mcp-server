package yunxiao

import (
	"context"
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
