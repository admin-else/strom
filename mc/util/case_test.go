package util_test

import (
	"testing"

	"git.anygate.cloud/anygatecloud/strom/mc/util"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "HelloWorld"},
		{"test_case", "TestCase"},
		{"simple-test", "SimpleTest"},
		{"already.Camel", "AlreadyCamel"},
		{"UPPER_CASE", "UPPERCASE"},
		{"mixed_Case_string", "MixedCaseString"},
		{"", ""},
	}

	for _, test := range tests {
		result := util.CamelCase(test.input)
		if result != test.expected {
			t.Errorf("CamelCase(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}
