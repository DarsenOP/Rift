// Package version contains build version information.
package version

import (
	"fmt"
	"regexp"
)

// Version is the current version of Rift.
// Follows semantic versioning: major.minor.patch-dev|release
const Version = "0.1.0-dev"

// ValidateVersion checks if the version string follows semantic versioning format.
func ValidateVersion(v string) error {
	// Regex for semantic versioning: d.d.d-dev OR d.d.d-release
	pattern := `^(\d+)\.(\d+)\.(\d+)-(dev|release)$`
	matched, err := regexp.MatchString(pattern, v)
	if err != nil {
		return fmt.Errorf("invalid version regex pattern: %v", err)
	}
	if !matched {
		return fmt.Errorf("version %s does not follow semantic versioning format: d.d.d-dev or d.d.d-release", v)
	}
	return nil
}
