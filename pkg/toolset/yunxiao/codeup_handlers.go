package yunxiao

import (
	"context"
	"net/url"
	"strings"
)

func handleListRepositories(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "search")
	setOptionalBool(query, params, "archived")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetRepository(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID)
	return c.GetJSON(ctx, path, nil)
}

func handleListBranches(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "search")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/branches"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetBranch(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	branchName, err := requiredString(params, "branchName")
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/branches/" + encodePathValue(branchName)
	return c.GetJSON(ctx, path, nil)
}

func handleListFiles(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalString(query, params, "path")
	setOptionalString(query, params, "ref")
	setOptionalString(query, params, "type")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/files/tree"
	return c.GetJSON(ctx, path, query)
}

func handleGetFileBlobs(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	filePath, err := requiredString(params, "filePath")
	if err != nil {
		return "", err
	}
	ref, err := requiredString(params, "ref")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("ref", ref)

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/files/" + encodeFilePath(filePath)
	return c.GetJSON(ctx, path, query)
}

func handleListCommits(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	refName, err := requiredString(params, "refName")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("refName", refName)
	setOptionalString(query, params, "since")
	setOptionalString(query, params, "until")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "path")
	setOptionalString(query, params, "search")
	setOptionalBool(query, params, "showSignature")
	setOptionalString(query, params, "committerIds")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/commits"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetCommit(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	sha, err := requiredString(params, "sha")
	if err != nil {
		return "", err
	}

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/commits/" + url.PathEscape(sha)
	return c.GetJSON(ctx, path, nil)
}

func handleCompare(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, err := requiredOrganizationAndRepository(params)
	if err != nil {
		return "", err
	}
	from, err := requiredString(params, "from")
	if err != nil {
		return "", err
	}
	to, err := requiredString(params, "to")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("from", from)
	query.Set("to", to)
	setOptionalString(query, params, "sourceType")
	setOptionalString(query, params, "targetType")
	setOptionalString(query, params, "straight")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/compares"
	return c.GetJSON(ctx, path, query)
}

func handleListChangeRequests(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "projectIds")
	setOptionalString(query, params, "authorIds")
	setOptionalString(query, params, "reviewerIds")
	setOptionalString(query, params, "state")
	setOptionalString(query, params, "search")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalString(query, params, "createdBefore")
	setOptionalString(query, params, "createdAfter")

	path := "/codeup/organizations/" + url.PathEscape(organizationID) + "/changeRequests"
	return c.GetJSONWithMetadata(ctx, path, query)
}

func handleGetChangeRequest(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}

	path := changeRequestPath(organizationID, repositoryID, localID)
	return c.GetJSON(ctx, path, nil)
}

func handleListChangeRequestPatchSets(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/diffs/patches"
	return c.GetJSON(ctx, path, nil)
}

func handleGetChangeRequestTree(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}
	fromPatchSetID, err := requiredString(params, "fromPatchSetId")
	if err != nil {
		return "", err
	}
	toPatchSetID, err := requiredString(params, "toPatchSetId")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("fromPatchSetId", fromPatchSetID)
	query.Set("toPatchSetId", toPatchSetID)

	path := changeRequestPath(organizationID, repositoryID, localID) + "/diffs/changeTree"
	return c.GetJSON(ctx, path, query)
}

func handleListChangeRequestComments(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}

	body := map[string]any{
		"comment_type": optionalStringDefault(params, "commentType", "GLOBAL_COMMENT"),
		"state":        optionalStringDefault(params, "state", "OPENED"),
		"resolved":     optionalBoolDefault(params, "resolved", false),
	}
	if patchSetBizIDs, _ := params["patchSetBizIds"].(string); patchSetBizIDs != "" {
		body["patchset_biz_id_list"] = strings.Join(splitCSV(patchSetBizIDs), ",")
	}
	if filePath, _ := params["filePath"].(string); filePath != "" {
		body["file_path"] = filePath
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/comments/list"
	return c.PostJSONWithMetadata(ctx, path, body)
}

func handleGetChangeRequestComment(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, repositoryID, localID, err := requiredOrganizationRepositoryAndLocalID(params)
	if err != nil {
		return "", err
	}
	commentBizID, err := requiredString(params, "commentBizId")
	if err != nil {
		return "", err
	}

	path := changeRequestPath(organizationID, repositoryID, localID) + "/comments/" + url.PathEscape(commentBizID)
	return c.PostJSONWithMetadata(ctx, path, nil)
}

func changeRequestPath(organizationID, repositoryID, localID string) string {
	return "/codeup/organizations/" + url.PathEscape(organizationID) + "/repositories/" + EncodeRepositoryID(repositoryID) + "/changeRequests/" + url.PathEscape(localID)
}
