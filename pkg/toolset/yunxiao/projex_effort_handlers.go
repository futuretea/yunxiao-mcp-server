package yunxiao

import (
	"context"
	"net/url"
)

func handleListCurrentUserEffortRecords(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, err := requiredString(params, "organizationId")
	if err != nil {
		return "", err
	}
	startDate, err := requiredString(params, "startDate")
	if err != nil {
		return "", err
	}
	endDate, err := requiredString(params, "endDate")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("startDate", startDate)
	query.Set("endDate", endDate)

	return c.GetJSON(ctx, projexOrganizationPath(organizationID)+"/effortRecords", query)
}

func handleListEffortRecords(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "workitemId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexWorkitemPath(organizationID, workItemID) + "/effortRecords"
	return c.GetJSON(ctx, path, nil)
}

func handleListEstimatedEfforts(ctx context.Context, client any, params map[string]any) (string, error) {
	organizationID, workItemID, err := requiredOrganizationAndNamedID(params, "workitemId")
	if err != nil {
		return "", err
	}
	c, err := getClient(client)
	if err != nil {
		return "", err
	}

	path := projexWorkitemPath(organizationID, workItemID) + "/estimatedEfforts"
	return c.GetJSON(ctx, path, nil)
}

func projexWorkitemPath(organizationID, workItemID string) string {
	return projexOrganizationPath(organizationID) + "/workitems/" + url.PathEscape(workItemID)
}
