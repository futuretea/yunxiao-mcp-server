package yunxiao

import (
	"context"
	"net/url"
)

func handleListKnowledgeBases(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}

	return c.GetJSONWithMetadata(ctx, lingmaOrganizationPath(organizationID)+"/knowledgeBases", lingmaKnowledgeBaseListQuery(params))
}

func handleListKbFiles(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, kbID, err := requiredOrganizationAndNamedID(params, "kbId")
	if err != nil {
		return "", err
	}

	path := lingmaKnowledgeBasePath(organizationID, kbID) + "/files"
	return c.GetJSONWithMetadata(ctx, path, lingmaKnowledgeBaseChildListQuery(params))
}

func handleListKbMembers(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, kbID, err := requiredOrganizationAndNamedID(params, "kbId")
	if err != nil {
		return "", err
	}

	path := lingmaKnowledgeBasePath(organizationID, kbID) + "/members"
	return c.GetJSONWithMetadata(ctx, path, lingmaKnowledgeBaseChildListQuery(params))
}

func lingmaKnowledgeBaseListQuery(params map[string]any) url.Values {
	query := lingmaKnowledgeBaseChildListQuery(params)
	setOptionalString(query, params, "sceneType")
	setOptionalString(query, params, "userId")
	return query
}

func lingmaKnowledgeBaseChildListQuery(params map[string]any) url.Values {
	query := url.Values{}
	setOptionalString(query, params, "query")
	setOptionalString(query, params, "orderBy")
	setOptionalString(query, params, "sort")
	setOptionalInt(query, params, "page")
	setOptionalInt(query, params, "perPage")
	return query
}

func lingmaKnowledgeBasePath(organizationID, kbID string) string {
	return lingmaOrganizationPath(organizationID) + "/knowledgeBases/" + encodePathValue(kbID)
}
