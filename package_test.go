package pkgen

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackagesQuery(t *testing.T) {
	tests := map[string]struct {
		config              PackagesQueryConfig
		expectedPackageName string
		expectError         bool
	}{
		"query current package": {
			config: PackagesQueryConfig{
				IncludeTests: false,
				Patterns:     []string{"."},
			},
			expectedPackageName: "pkgen",
			expectError:         false,
		},
		"query with tests": {
			config: PackagesQueryConfig{
				IncludeTests: true,
				Patterns:     []string{"."},
			},
			expectedPackageName: "pkgen",
			expectError:         false,
		},
		"invalid pattern": {
			config: PackagesQueryConfig{
				Patterns: []string{"./nonexistent/package/that/does/not/exist"},
			},
			expectError: false, // packages.Load doesn't error on non-existent patterns, just returns empty
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Packages{}
			pkgs, err := p.Query(context.Background(), tc.config)

			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tc.expectedPackageName != "" {
				require.NotEmpty(t, pkgs, "expected at least one package")
				assert.Equal(t, tc.expectedPackageName, pkgs[0].Name)
			}
		})
	}
}
