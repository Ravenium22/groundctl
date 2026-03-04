package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// CheckConstraint checks whether a raw version string satisfies a constraint string.
// Supported constraint formats:
//   - Exact:  "1.2.3"
//   - Range:  ">=1.2.0", ">=1.2.0, <2.0.0"
//   - Caret:  "^1.2.3" (>=1.2.3, <2.0.0)
//   - Tilde:  "~1.2.3" (>=1.2.3, <1.3.0)
//   - Star:   "*" (any version)
//
// Returns (satisfied, error). If parsing fails, returns an error.
func CheckConstraint(constraintStr, versionStr string) (bool, error) {
	versionStr = normalizeVersion(versionStr)
	constraintStr = normalizeConstraint(constraintStr)

	if constraintStr == "" || constraintStr == "*" {
		return true, nil
	}

	ver, err := semver.NewVersion(versionStr)
	if err != nil {
		return false, fmt.Errorf("invalid version %q: %w", versionStr, err)
	}

	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return false, fmt.Errorf("invalid constraint %q: %w", constraintStr, err)
	}

	return constraint.Check(ver), nil
}

var leadingV = regexp.MustCompile(`^v`)

// normalizeVersion strips leading 'v' and pads to at least x.y.z.
func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = leadingV.ReplaceAllString(v, "")

	parts := strings.SplitN(v, ".", 3)
	switch len(parts) {
	case 1:
		return parts[0] + ".0.0"
	case 2:
		return parts[0] + "." + parts[1] + ".0"
	default:
		return v
	}
}

// normalizeConstraint handles our custom constraint syntax.
func normalizeConstraint(c string) string {
	c = strings.TrimSpace(c)
	if c == "" || c == "*" {
		return c
	}

	// Pad bare versions in constraints (e.g. "^3.10" -> "^3.10.0")
	// Split on comma for compound constraints
	parts := strings.Split(c, ",")
	for i, part := range parts {
		parts[i] = padConstraintPart(strings.TrimSpace(part))
	}
	return strings.Join(parts, ", ")
}

var constraintPrefixRe = regexp.MustCompile(`^([~^>=<!]+)?\s*(.+)$`)

func padConstraintPart(part string) string {
	matches := constraintPrefixRe.FindStringSubmatch(part)
	if matches == nil {
		return part
	}

	prefix := matches[1]
	ver := matches[2]
	ver = leadingV.ReplaceAllString(ver, "")

	segments := strings.SplitN(ver, ".", 3)
	switch len(segments) {
	case 1:
		ver = segments[0] + ".0.0"
	case 2:
		ver = segments[0] + "." + segments[1] + ".0"
	}

	return prefix + ver
}
