package yunxiao

import (
	"context"
	"net/url"
)

func handleListAuditLogs(ctx context.Context, client any, params map[string]any) (string, error) {
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	actionTimeStart, err := requiredString(params, "actionTimeStart")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("actionTimeStart", actionTimeStart)
	setOptionalString(query, params, "userIds")
	setOptionalString(query, params, "apps")
	setOptionalString(query, params, "actionTimeEnd")
	setOptionalInt(query, params, "perPage")
	setOptionalString(query, params, "nextToken")

	return c.GetJSONWithMetadata(ctx, organizationPath(organizationID)+"/auditLogs", query)
}
