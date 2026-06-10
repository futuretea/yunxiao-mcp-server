package yunxiao

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// PrettyJSON returns an indented JSON string, or the raw body when it is not valid JSON.
func PrettyJSON(raw json.RawMessage) string {
	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return string(raw)
	}
	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return string(raw)
	}
	return string(formatted)
}

// PrettyResponseJSON returns an indented JSON payload that includes response metadata.
func PrettyResponseJSON(resp *Response) string {
	var data any
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		data = string(resp.Body)
	}

	payload := map[string]any{"data": data}
	if resp.Pagination != nil {
		payload["pagination"] = resp.Pagination
	}
	if resp.NextToken != "" {
		payload["nextToken"] = resp.NextToken
	}
	if resp.RequestID != "" {
		payload["requestId"] = resp.RequestID
	}
	formatted, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return PrettyJSON(resp.Body)
	}
	return string(formatted)
}

func parsePagination(header http.Header) *Pagination {
	pagination := &Pagination{
		Page:       parseHeaderInt(header, "x-page"),
		PerPage:    parseHeaderInt(header, "x-per-page"),
		PrevPage:   parseHeaderInt(header, "x-prev-page"),
		NextPage:   parseHeaderInt(header, "x-next-page"),
		Total:      parseHeaderInt(header, "x-total"),
		TotalPages: parseHeaderInt(header, "x-total-pages"),
	}
	if pagination.Page == 0 &&
		pagination.PerPage == 0 &&
		pagination.PrevPage == 0 &&
		pagination.NextPage == 0 &&
		pagination.Total == 0 &&
		pagination.TotalPages == 0 {
		return nil
	}
	return pagination
}

func parseHeaderInt(header http.Header, key string) int {
	value := header.Get(key)
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}
