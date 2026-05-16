package yunxiao

import (
	"context"
	"net/url"
	"strings"
)

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
