package yunxiao

import (
	"net/url"
	"strings"
)

func (c *Client) resolveURL(path string, query url.Values) string {
	u := *c.baseURL
	escapedPath := joinEscapedPath(c.baseURL.EscapedPath(), path)
	if decodedPath, err := url.PathUnescape(escapedPath); err == nil {
		u.Path = decodedPath
		u.RawPath = escapedPath
	} else {
		// Invalid escape sequences are rare; fall back to a plain concatenation
		// and let the server decide how to handle the malformed path.
		u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
		u.RawPath = ""
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func joinEscapedPath(basePath, path string) string {
	basePath = strings.TrimRight(basePath, "/")
	path = strings.TrimLeft(path, "/")
	if basePath == "" {
		return "/" + path
	}
	if path == "" {
		return basePath
	}
	return basePath + "/" + path
}

// EncodeRepositoryID encodes a CodeUp repository ID or full path for path usage.
func EncodeRepositoryID(repositoryID string) string {
	repositoryID = strings.TrimSpace(repositoryID)
	if repositoryID == "" {
		return ""
	}
	if strings.Contains(repositoryID, "%2F") || strings.Contains(repositoryID, "%2f") {
		return repositoryID
	}
	if !strings.Contains(repositoryID, "/") {
		return url.PathEscape(repositoryID)
	}

	parts := strings.SplitN(repositoryID, "/", 2)
	return url.PathEscape(parts[0]) + "%2F" + url.PathEscape(parts[1])
}

// EncodePathValue escapes a single Yunxiao path value while preserving existing %2F separators.
func EncodePathValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "%2F") || strings.Contains(value, "%2f") {
		return value
	}
	return url.PathEscape(value)
}

// EncodeFilePath escapes a repository file path for CodeUp file endpoints.
func EncodeFilePath(filePath string) string {
	return EncodePathValue(strings.TrimPrefix(strings.TrimSpace(filePath), "/"))
}
