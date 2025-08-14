// Package pkgen provides utilities for Go package discovery and code generation.
//
// This package offers functionality to query Go packages using configurable
// patterns and generate code files within those packages using templates.
// It supports both standalone usage and integration with go:generate directives.
//
// Key features:
//   - Package discovery using golang.org/x/tools/go/packages
//   - Template-based code generation
//   - YAML configuration support
//   - go:generate integration
package pkgen
