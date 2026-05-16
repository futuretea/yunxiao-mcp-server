package yunxiao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type projectOverviewSection struct {
	flag  string
	name  string
	path  string
	query url.Values
}

func projectOverviewSections(projectPath string, params map[string]any) []projectOverviewSection {
	return []projectOverviewSection{
		{flag: "includeMembers", name: "members", path: projectPath + "/members"},
		{flag: "includeSprints", name: "sprints", path: projectPath + "/sprints", query: projectOverviewListQuery(params, true)},
		{flag: "includeMilestones", name: "milestones", path: projectPath + "/milestones", query: projectOverviewListQuery(params, true)},
		{flag: "includeVersions", name: "versions", path: projectPath + "/versions", query: projectOverviewListQuery(params, true)},
		{flag: "includeLabels", name: "labels", path: projectPath + "/labels", query: projectOverviewListQuery(params, false)},
	}
}

func addProjectOverviewSection(ctx context.Context, c *Client, overview map[string]any, params map[string]any, section projectOverviewSection) error {
	if !optionalBoolDefault(params, section.flag, true) {
		return nil
	}
	payload, err := getProjectOverviewSection(ctx, c, section.name, section.path, section.query)
	if err != nil {
		return err
	}
	overview[section.name] = payload
	return nil
}

func getProjectOverviewSection(ctx context.Context, c *Client, name, path string, query url.Values) (any, error) {
	resp, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return responsePayload(resp), nil
}

func projectOverviewListQuery(params map[string]any, withStatus bool) url.Values {
	query := url.Values{}
	query.Set("page", strconv.Itoa(optionalIntDefault(params, "page", 1)))
	query.Set("perPage", strconv.Itoa(optionalIntDefault(params, "perPage", 20)))
	if withStatus && optionalBoolDefault(params, "activeOnly", true) {
		status := optionalStringDefault(params, "status", "TODO,DOING")
		if status != "" {
			query.Set("status", status)
		}
	}
	return query
}

func projectOverviewFilters(params map[string]any) map[string]any {
	return map[string]any{
		"activeOnly": optionalBoolDefault(params, "activeOnly", true),
		"status":     optionalStringDefault(params, "status", "TODO,DOING"),
		"page":       optionalIntDefault(params, "page", 1),
		"perPage":    optionalIntDefault(params, "perPage", 20),
	}
}

func addProjectWorkitemContextBaseSections(ctx context.Context, c *Client, payload map[string]any, params map[string]any, projectPath, category string) error {
	typeQuery := url.Values{}
	typeQuery.Set("category", category)
	types, err := getProjectOverviewSection(ctx, c, "workItemTypes", projectPath+"/workitemTypes", typeQuery)
	if err != nil {
		return err
	}
	payload["workItemTypes"] = types

	if optionalBoolDefault(params, "includeMembers", true) {
		members, err := getProjectOverviewSection(ctx, c, "members", projectPath+"/members", nil)
		if err != nil {
			return err
		}
		payload["members"] = members
	}
	if optionalBoolDefault(params, "includeLabels", true) {
		labels, err := getProjectOverviewSection(ctx, c, "labels", projectPath+"/labels", projectOverviewListQuery(params, false))
		if err != nil {
			return err
		}
		payload["labels"] = labels
	}
	return nil
}

func addProjectWorkitemTypeContext(ctx context.Context, c *Client, payload map[string]any, params map[string]any, projectPath string) error {
	workItemTypeID, _ := params["workItemTypeId"].(string)
	if strings.TrimSpace(workItemTypeID) == "" {
		return nil
	}
	typePath := projectPath + "/workitemTypes/" + url.PathEscape(strings.TrimSpace(workItemTypeID))
	if optionalBoolDefault(params, "includeFields", true) {
		fields, err := getProjectOverviewSection(ctx, c, "fields", typePath+"/fields", nil)
		if err != nil {
			return err
		}
		payload["fields"] = fields
	}
	if optionalBoolDefault(params, "includeWorkflow", true) {
		workflow, err := getProjectOverviewSection(ctx, c, "workflow", typePath+"/workflows", nil)
		if err != nil {
			return err
		}
		payload["workflow"] = workflow
	}
	return nil
}

type workitemDetailSection struct {
	flag  string
	name  string
	path  string
	query url.Values
}

func workitemDetailSections(workitemPath string, params map[string]any) []workitemDetailSection {
	sections := []workitemDetailSection{
		{flag: "includeActivities", name: "activities", path: workitemPath + "/activities"},
		{flag: "includeAttachments", name: "attachments", path: workitemPath + "/attachments"},
	}

	if optionalBoolDefault(params, "includeComments", true) {
		query := url.Values{}
		query.Set("page", strconv.Itoa(optionalIntDefault(params, "page", 1)))
		query.Set("perPage", strconv.Itoa(optionalIntDefault(params, "perPage", 20)))
		sections = append(sections, workitemDetailSection{flag: "includeComments", name: "comments", path: workitemPath + "/comments", query: query})
	}

	if optionalBoolDefault(params, "includeRelations", true) {
		relationTypes := splitCSV(optionalStringDefault(params, "relationTypes", "ASSOCIATED,SUB"))
		for _, rt := range relationTypes {
			query := url.Values{}
			query.Set("relationType", rt)
			sections = append(sections, workitemDetailSection{
				flag:  "includeRelations",
				name:  "relations_" + strings.ToLower(rt),
				path:  workitemPath + "/relationRecords",
				query: query,
			})
		}
	}

	return sections
}

func addWorkitemDetailSection(ctx context.Context, c *Client, detail map[string]any, params map[string]any, section workitemDetailSection) error {
	if !optionalBoolDefault(params, section.flag, true) {
		return nil
	}
	payload, err := getProjectOverviewSection(ctx, c, section.name, section.path, section.query)
	if err != nil {
		return err
	}
	detail[section.name] = payload
	return nil
}

func workitemDetailFilters(params map[string]any) map[string]any {
	return map[string]any{
		"includeActivities":  optionalBoolDefault(params, "includeActivities", true),
		"includeRelations":   optionalBoolDefault(params, "includeRelations", true),
		"relationTypes":      optionalStringDefault(params, "relationTypes", "ASSOCIATED,SUB"),
		"includeAttachments": optionalBoolDefault(params, "includeAttachments", true),
		"includeComments":    optionalBoolDefault(params, "includeComments", true),
		"page":               optionalIntDefault(params, "page", 1),
		"perPage":            optionalIntDefault(params, "perPage", 20),
	}
}
