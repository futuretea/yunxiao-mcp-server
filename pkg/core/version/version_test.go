package version

import (
	"strings"
	"testing"
)

func TestGetVersionInfoIncludesBinaryName(t *testing.T) {
	info := GetVersionInfo()
	if !strings.Contains(info, BinaryName) {
		t.Fatalf("GetVersionInfo() = %q, missing binary name", info)
	}
	if !strings.Contains(info, "version=") {
		t.Fatalf("GetVersionInfo() = %q, missing version", info)
	}
	if !strings.Contains(info, "commit=") {
		t.Fatalf("GetVersionInfo() = %q, missing commit", info)
	}
	if !strings.Contains(info, "date=") {
		t.Fatalf("GetVersionInfo() = %q, missing date", info)
	}
}
