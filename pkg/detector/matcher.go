package detector

import (
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version.
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string into a Version struct.
func ParseVersion(v string) (*Version, error) {
	// Remove 'v' prefix if present
	v = strings.TrimPrefix(v, "v")

	// Try to parse as semantic version first
	parts := strings.Split(v, ".")
	if len(parts) >= 1 {
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		minor := 0
		if len(parts) >= 2 {
			minor, _ = strconv.Atoi(parts[1])
		}

		patch := 0
		if len(parts) >= 3 {
			// Handle versions like "1.2.3-rc1" by taking just the numeric part
			patchStr := strings.Split(parts[2], "-")[0]
			patch, _ = strconv.Atoi(patchStr)
		}

		return &Version{Major: major, Minor: minor, Patch: patch}, nil
	}

	return nil, ErrInvalidVersion
}

// CompareVersions compares two version strings.
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
func CompareVersions(v1, v2 string) (int, error) {
	ver1, err := ParseVersion(v1)
	if err != nil {
		return 0, err
	}

	ver2, err := ParseVersion(v2)
	if err != nil {
		return 0, err
	}

	return CompareVersionStructs(ver1, ver2), nil
}

// CompareVersionStructs compares two Version structs.
func CompareVersionStructs(v1, v2 *Version) int {
	if v1.Major != v2.Major {
		if v1.Major < v2.Major {
			return -1
		}
		return 1
	}

	if v1.Minor != v2.Minor {
		if v1.Minor < v2.Minor {
			return -1
		}
		return 1
	}

	if v1.Patch != v2.Patch {
		if v1.Patch < v2.Patch {
			return -1
		}
		return 1
	}

	return 0
}

// VersionGreaterThan returns true if v1 > v2.
func VersionGreaterThan(v1, v2 string) bool {
	cmp, _ := CompareVersions(v1, v2)
	return cmp > 0
}

// VersionLessThan returns true if v1 < v2.
func VersionLessThan(v1, v2 string) bool {
	cmp, _ := CompareVersions(v1, v2)
	return cmp < 0
}

// VersionEqual returns true if v1 == v2.
func VersionEqual(v1, v2 string) bool {
	cmp, _ := CompareVersions(v1, v2)
	return cmp == 0
}

// VersionAtLeast returns true if v1 >= v2.
func VersionAtLeast(v1, v2 string) bool {
	cmp, _ := CompareVersions(v1, v2)
	return cmp >= 0
}

// VersionAtMost returns true if v1 <= v2.
func VersionAtMost(v1, v2 string) bool {
	cmp, _ := CompareVersions(v1, v2)
	return cmp <= 0
}

// MatchVersionPattern checks if a version string matches a given pattern.
// Supports patterns like "1.x", "1.2.x", ">=1.0.0", "<2.0.0", etc.
func MatchVersionPattern(version, pattern string) (bool, error) {
	// Handle range patterns
	if strings.HasPrefix(pattern, ">=") {
		return VersionAtLeast(version, strings.TrimPrefix(pattern, ">=")), nil
	}
	if strings.HasPrefix(pattern, "<=") {
		return VersionAtMost(version, strings.TrimPrefix(pattern, "<=")), nil
	}
	if strings.HasPrefix(pattern, ">") {
		return VersionGreaterThan(version, strings.TrimPrefix(pattern, ">")), nil
	}
	if strings.HasPrefix(pattern, "<") {
		return VersionLessThan(version, strings.TrimPrefix(pattern, "<")), nil
	}

	// Handle x/y/z wildcards
	if strings.Contains(pattern, "x") {
		return matchVersionWildcard(version, pattern), nil
	}

	// Exact match
	return VersionEqual(version, pattern), nil
}

// matchVersionWildcard matches version with x wildcards (e.g., "1.x", "1.2.x").
func matchVersionWildcard(version, pattern string) bool {
	versionParts := strings.Split(version, ".")
	patternParts := strings.Split(pattern, ".")

	for i, p := range patternParts {
		if p == "x" || p == "*" {
			// Wildcard matches anything
			continue
		}
		if i >= len(versionParts) {
			// Pattern has more specificity than version
			return p == "0"
		}
		if versionParts[i] != p {
			return false
		}
	}

	return true
}

// MatchVersionRegex matches a version against a regex pattern.
func MatchVersionRegex(version, regexPattern string) bool {
	re := regexp.MustCompile(regexPattern)
	return re.MatchString(version)
}

// IsValidVersion checks if a string is a valid version format.
func IsValidVersion(v string) bool {
	_, err := ParseVersion(v)
	return err == nil
}

// String returns the version as a string.
func (v *Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

// ErrInvalidVersion is returned when a version string cannot be parsed.
var ErrInvalidVersion = &invalidVersionError{}

type invalidVersionError struct{}

func (e *invalidVersionError) Error() string {
	return "invalid version format"
}
