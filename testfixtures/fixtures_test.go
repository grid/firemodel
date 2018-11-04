package testfixtures

import (
	"gotest.tools/assert"
	"testing"
)
import firemodels "github.com/visor-tax/firemodel/testfixtures/firemodel/TestFiremodelFromSchema/go"

func TestRegexPath(t *testing.T) {
	for _, tt := range []struct {
		name string
		arg  string
		exp  bool
	}{
		{"empty", "", false},
		{"fully qualified", "/projects/some-project/databases/(default)/documents/users/123/test_models/abc", true},
		{"doc only", "/users/123/test_models/abc", true},
		{"fully qualified no leading slash", "projects/some-project/databases/(default)/documents/users/123/test_models/abc", true},
		{"doc only no leading slash", "users/123/test_models/abc", true},

		{"prefix match, not real match", "users/123", false},
		{"prefix match, not real match", "/users/123", false},
		{"prefix match, not real match", "/users/123/", false},
		{"prefix match, not real match", "users/123/", false},

		{"non match", "/othermodel/random", false},
		{"roundtrip", firemodels.TestModelPath("userid", "testmodelid"), true},
		{"empty ids", firemodels.TestModelPath("", ""), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, firemodels.TestModelRegexPath.MatchString(tt.arg), tt.exp)
		})
	}
}
