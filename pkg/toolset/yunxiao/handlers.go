package yunxiao

import (
	"bytes"
	"encoding/json"
)

func marshalPretty(value any) (string, error) {
	formatted, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

func responsePayload(resp *Response) any {
	var data any
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		data = string(resp.Body)
	}

	if resp.Pagination == nil && resp.NextToken == "" && resp.RequestID == "" {
		return data
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
	return payload
}

func responsePayloadPreservingNumbers(resp *Response) any {
	data := decodeJSONPreservingNumbers(resp.Body)
	if resp.Pagination == nil && resp.NextToken == "" && resp.RequestID == "" {
		return data
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
	return payload
}

func decodeJSONPreservingNumbers(raw []byte) any {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var data any
	if err := decoder.Decode(&data); err != nil {
		return string(raw)
	}
	return data
}

func prettyJSONPreservingNumbers(raw []byte) string {
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, raw, "", "  "); err != nil {
		return string(raw)
	}
	return formatted.String()
}
