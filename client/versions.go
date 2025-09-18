package client

import "fmt"

const MaxSupportedVersion int = 65
const MinSupportedVersion int = 40

const DefaultVersion = MaxSupportedVersion

func validateVersion(version int) bool {
	return version <= MaxSupportedVersion &&
		version >= MinSupportedVersion
}

func toVersionString(version int) string {
	return fmt.Sprintf(
		"%d.0",
		version,
	)
}
