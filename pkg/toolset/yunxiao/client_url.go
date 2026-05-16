package yunxiao

import (
	"net/url"
	"strings"
)

func (c *Client) resolveURL(path string, query url.Values) string {
	u := *c.baseURL
	escapedPath := joinEscapedPath(c.baseURL.EscapedPath(), path)
	decodedPath, err := url.PathUnescape(escapedPath)
	if err == nil {
		u.Path = decodedPath
		u.RawPath = escapedPath
	} else {
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

func encodePathValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "%2F") || strings.Contains(value, "%2f") {
		return value
	}
	return url.PathEscape(value)
}

func encodeFilePath(filePath string) string {
	return encodePathValue(strings.TrimPrefix(strings.TrimSpace(filePath), "/"))
}
