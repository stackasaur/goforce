package client

import "testing"

func TestValidateVersion(t *testing.T) {
	validVersion := MaxSupportedVersion - 1

	expected := true
	actual := validateVersion(validVersion)

	if expected != actual {
		t.Fatalf(
			"expected %v, actual %v",
			expected,
			actual,
		)
	}

	invalidVersion := MaxSupportedVersion + 1

	expected = false
	actual = validateVersion(invalidVersion)

	if expected != actual {
		t.Fatalf(
			"expected %v, actual %v",
			expected,
			actual,
		)
	}
}
func TestToVersionString(t *testing.T) {
	version := 60

	expected := "60.0"
	actual := toVersionString(version)

	if expected != actual {
		t.Fatalf(
			"expected %v, actual %v",
			expected,
			actual,
		)
	}
}
