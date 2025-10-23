package pkgen

import (
	"context"
	"testing"

	"github.com/ifnotnil/x/tst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackagesQuery(t *testing.T) {
	tests := map[string]struct {
		config              PackagesQueryConfig
		expectedPackageName string
		errorAsserter       tst.ErrorAssertionFunc
	}{
		"query current package": {
			config: PackagesQueryConfig{
				IncludeTests: false,
				Patterns:     []string{"."},
			},
			expectedPackageName: "pkgen",
			errorAsserter:       tst.NoError(),
		},
		"query with tests": {
			config: PackagesQueryConfig{
				IncludeTests: true,
				Patterns:     []string{"."},
			},
			expectedPackageName: "pkgen",
			errorAsserter:       tst.NoError(),
		},
		"invalid pattern": {
			config: PackagesQueryConfig{
				Patterns: []string{"./nonexistent/package/that/does/not/exist"},
			},
			expectedPackageName: "",
			errorAsserter:       tst.NoError(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Packages{}
			pkgs, err := p.Query(context.Background(), tc.config)

			tc.errorAsserter(t, err)

			require.NoError(t, err)

			if tc.expectedPackageName != "" {
				require.NotEmpty(t, pkgs, "expected at least one package")
				assert.Equal(t, tc.expectedPackageName, pkgs[0].Name)
			}
		})
	}
}
