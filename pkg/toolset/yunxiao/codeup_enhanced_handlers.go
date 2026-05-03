package yunxiao

import (
	"context"
	"net/url"
	"strconv"
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

	repoPath := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID)

	repository, err := getProjectOverviewSection(ctx, c, "repository", repoPath, nil)
	if err != nil {
		return "", err
	}

	overview := map[string]any{
		"repository": repository,
		"filters":    repositoryOverviewFilters(params),
	}

	if optionalBoolDefault(params, "includeBranches", true) {
		branches, err := getProjectOverviewSection(ctx, c, "branches", repoPath+"/branches", repositoryLimitQuery(params, "branchLimit", 5))
		if err != nil {
			return "", err
		}
		overview["branches"] = branches
	}

	refName := optionalStringDefault(params, "refName", "")
	if refName == "" {
		if repoMap, ok := repository.(map[string]any); ok {
			if db, ok := repoMap["defaultBranch"].(string); ok && db != "" {
				refName = db
			}
		}
	}

	if optionalBoolDefault(params, "includeCommits", true) && refName != "" {
		commitQuery := repositoryLimitQuery(params, "commitLimit", 5)
		commitQuery.Set("refName", refName)
		commits, err := getProjectOverviewSection(ctx, c, "commits", repoPath+"/commits", commitQuery)
		if err != nil {
			return "", err
		}
		overview["commits"] = commits
	}

	if optionalBoolDefault(params, "includeMergeRequests", true) {
		mrQuery := repositoryLimitQuery(params, "mrLimit", 5)
		mrQuery.Set("state", optionalStringDefault(params, "mrState", "opened"))
		mrQuery.Add("repositoryIds", repositoryID)
		mrs, err := getProjectOverviewSection(ctx, c, "mergeRequests", codeupOrganizationPath(organizationID)+"/mergeRequests", mrQuery)
		if err != nil {
			return "", err
		}
		overview["mergeRequests"] = mrs
	}

	return marshalPretty(overview)
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

func repositoryLimitQuery(params map[string]any, limitKey string, defaultLimit int) url.Values {
	query := url.Values{}
	query.Set("page", "1")
	query.Set("perPage", strconv.Itoa(optionalIntDefault(params, limitKey, defaultLimit)))
	return query
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
