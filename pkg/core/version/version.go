package version

import "fmt"

const BinaryName = "yunxiao"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// GetVersionInfo returns human-readable build metadata.
func GetVersionInfo() string {
	return GetVersionInfoFor(BinaryName)
}

// GetVersionInfoFor returns human-readable build metadata for the named binary.
func GetVersionInfoFor(binaryName string) string {
	return fmt.Sprintf("%s version=%s commit=%s date=%s", binaryName, Version, Commit, Date)
}
