package dpkg

import (
	"regexp"
	"strconv"
	"strings"
)

// CompareVersions compares Debian package versions based on https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-version
func (dp *Dpkg) CompareVersions(v1, v2 string) int {
	epoch1, version1, revision1 := splitVersion(v1)
	epoch2, version2, revision2 := splitVersion(v2)

	// Compare epochs
	if epoch1 != epoch2 {
		return epoch1 - epoch2
	}

	// Compare main versions
	if cmp := compareDebianVersion(version1, version2); cmp != 0 {
		return cmp
	}

	// Compare revisions
	return compareDebianVersion(revision1, revision2)
}

// splitVersion splits a version into epoch, main version, and revision
func splitVersion(version string) (epoch int, mainVersion, revision string) {
	// Regex to capture epoch, main version, and revision
	re := regexp.MustCompile(`^(?:(\d+):)?([^-]+)(?:-(.+))?$`)
	matches := re.FindStringSubmatch(version)

	// Extract epoch (if exists)
	if matches[1] != "" {
		epoch, _ = strconv.Atoi(matches[1])
	} else {
		epoch = 0 // Default epoch is 0
	}

	// Extract main version and revision
	mainVersion = matches[2]
	revision = matches[3]

	return
}

// compareDebianVersion compares parts of the version (main version or revision)
func compareDebianVersion(v1, v2 string) int {
	// Split versions into parts (numbers and strings)
	parts1 := splitVersionParts(v1)
	parts2 := splitVersionParts(v2)

	// Compare each part
	for i := 0; i < len(parts1) || i < len(parts2); i++ {
		part1 := getPart(parts1, i)
		part2 := getPart(parts2, i)

		if part1 == part2 {
			continue
		}

		// Compare numerically if both parts are numbers
		num1, err1 := strconv.Atoi(part1)
		num2, err2 := strconv.Atoi(part2)
		if err1 == nil && err2 == nil {
			if num1 < num2 {
				return -1
			} else if num1 > num2 {
				return 1
			}
		} else {
			// Compare lexicographically
			if part1 < part2 {
				return -1
			} else if part1 > part2 {
				return 1
			}
		}
	}

	return 0
}

// splitVersionParts splits a version into parts (numbers and strings)
func splitVersionParts(version string) []string {
	var parts []string
	var current strings.Builder

	for _, char := range version {
		if char >= '0' && char <= '9' {
			// If it's a digit, add to the current number
			current.WriteRune(char)
		} else {
			// If it's not a digit, finalize the current number (if any)
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			// Add the non-numeric character as a separate part
			parts = append(parts, string(char))
		}
	}

	// Add the last number (if any)
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// getPart retrieves a part from a slice, returning "" if the index is invalid
func getPart(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}
