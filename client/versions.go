package client

import "regexp"

const DefaultVersion string = "65.0"

var VersionRegexp *regexp.Regexp = regexp.MustCompile(`^\d+.\d$`)

func ValidateVersion(version string) bool {
	return VersionRegexp.MatchString(version)
}
