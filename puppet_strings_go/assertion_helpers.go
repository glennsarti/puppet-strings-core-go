package puppet_strings_go

import (
	"strings"
	"testing"
)

func assertInteger(t *testing.T, fs string, prefix string, expect int, actual int) (bool) {
	if actual != expect { t.Errorf(fs, prefix, expect, actual) }
	return actual == expect
}

func assertString(t *testing.T, fs string, prefix string, expect string, actual string) (bool) {
	if actual != expect { t.Errorf(fs, prefix, expect, actual) }
	return actual == expect
}

func assertStringArray(t *testing.T, prefix string, actual []string, expect []string) (bool) {
	// Check for nil-ness
	if expect == nil && actual == nil { return true}
	if expect == nil && actual != nil {
		t.Errorf("%s: Expected list to be nil but got a list", prefix)
		return false
	}
	if expect != nil && actual == nil {
		t.Errorf("%s: Expected a list but got nil", prefix)
		return false
	}
	if len(actual) != len(expect) {
		t.Errorf(
			"%s: Expected list to be [%s], but got [%s]",
			prefix,
			strings.Join(expect, ", "),
			strings.Join(actual, ", "),
		)
		return false
	} else {
		pass := true
		for i, item := range expect {
			if actual[i] != item {
				t.Errorf("%s: Expected list item %d to be '%s' but got '%s'", prefix, i, item, actual[i])
				pass = false
			}
		}
		return pass
	}
}
