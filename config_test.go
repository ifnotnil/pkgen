package pkgen

import (
	"embed"
	"flag"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/ifnotnil/x/tst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

func TestPackagesQueryConfigFlags(t *testing.T) {
	tests := []struct {
		arguments []string
		expected  PackagesQueryConfig
	}{
		{
			arguments: []string{},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: nil, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-include_tests"},
			expected:  PackagesQueryConfig{IncludeTests: true, Env: nil, BuildFlags: nil, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"--include_tests"},
			expected:  PackagesQueryConfig{IncludeTests: true, Env: nil, BuildFlags: nil, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-env", "GOOS=linux"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: []string{"GOOS=linux"}, BuildFlags: nil, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"--env", "GOOS=linux"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: []string{"GOOS=linux"}, BuildFlags: nil, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-env", "GOOS=linux", "-env", "GOARCH=amd64"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: []string{"GOOS=linux", "GOARCH=amd64"}, BuildFlags: nil, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-build_flag", "-tags=debug"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: []string{"-tags=debug"}, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"--build_flag", "-tags=debug"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: []string{"-tags=debug"}, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-build_flag", "-tags=debug", "-build_flag", "-race"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: []string{"-tags=debug", "-race"}, Dir: "", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-dir", "/path/to/project"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: nil, Dir: "/path/to/project", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"--dir", "/path/to/project"},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: nil, Dir: "/path/to/project", Patterns: []string{"./..."}},
		},
		{
			arguments: []string{"-pattern", "./cmd/..."},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: nil, Dir: "", Patterns: []string{"./...", "./cmd/..."}},
		},
		{
			arguments: []string{"--pattern", "./cmd/..."},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: nil, Dir: "", Patterns: []string{"./...", "./cmd/..."}},
		},
		{
			arguments: []string{"-pattern", "./cmd/...", "-pattern", "./internal/..."},
			expected:  PackagesQueryConfig{IncludeTests: false, Env: nil, BuildFlags: nil, Dir: "", Patterns: []string{"./...", "./cmd/...", "./internal/..."}},
		},
		{
			arguments: []string{"-include_tests", "-env", "GOOS=linux", "-build_flag", "-race", "-dir", "/tmp", "-pattern", "./cmd/..."},
			expected:  PackagesQueryConfig{IncludeTests: true, Env: []string{"GOOS=linux"}, BuildFlags: []string{"-race"}, Dir: "/tmp", Patterns: []string{"./...", "./cmd/..."}},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fs := flag.NewFlagSet("tc", flag.ContinueOnError)
			c := DefaultConfig.PackagesQuery
			c.RegisterFlags(fs)
			err := fs.Parse(tc.arguments)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, c)
		})
	}
}

func TestGenerateConfigFlags(t *testing.T) {
	tests := []struct {
		arguments []string
		expected  GenerateConfig
	}{
		{
			arguments: []string{},
			expected:  GenerateConfig{OutputFile: defaultOutputNameTemplate, OutputFileMod: os.FileMode(0o644)},
		},
		{
			arguments: []string{"-output", "custom.go"},
			expected:  GenerateConfig{OutputFile: "custom.go", OutputFileMod: os.FileMode(0o644)},
		},
		{
			arguments: []string{"--output", "custom.go"},
			expected:  GenerateConfig{OutputFile: "custom.go", OutputFileMod: os.FileMode(0o644)},
		},
		{
			arguments: []string{"-mod", "0o755"},
			expected:  GenerateConfig{OutputFile: defaultOutputNameTemplate, OutputFileMod: os.FileMode(0o755)},
		},
		{
			arguments: []string{"--mod", "0o755"},
			expected:  GenerateConfig{OutputFile: defaultOutputNameTemplate, OutputFileMod: os.FileMode(0o755)},
		},
		{
			arguments: []string{"-mod", "0O644"},
			expected:  GenerateConfig{OutputFile: defaultOutputNameTemplate, OutputFileMod: os.FileMode(0o644)},
		},
		{
			arguments: []string{"-mod", "600"},
			expected:  GenerateConfig{OutputFile: defaultOutputNameTemplate, OutputFileMod: os.FileMode(0o600)},
		},
		{
			arguments: []string{"-output", "test.go", "-mod", "0o755"},
			expected:  GenerateConfig{OutputFile: "test.go", OutputFileMod: os.FileMode(0o755)},
		},
		{
			arguments: []string{"--output", "test.go", "--mod", "0o755"},
			expected:  GenerateConfig{OutputFile: "test.go", OutputFileMod: os.FileMode(0o755)},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fs := flag.NewFlagSet("tc", flag.ContinueOnError)
			c := DefaultConfig.Generate
			c.RegisterFlags(fs)
			err := fs.Parse(tc.arguments)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, c)
		})
	}
}

func TestTemplateConfigYAMLUnmarshal(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected TemplateConfigs
	}{
		"string": {
			input:    `"a single string"`,
			expected: TemplateConfigs{TemplateConfig{Name: "a single string", CustomTemplateFile: ""}},
		},
		"string array": {
			input:    `[ "abc", "def" ]`,
			expected: TemplateConfigs{TemplateConfig{Name: "abc", CustomTemplateFile: ""}, TemplateConfig{Name: "def", CustomTemplateFile: ""}},
		},
		"object array": {
			input: `- name: "abc"
- template_file: "/abc/def"`,
			expected: TemplateConfigs{TemplateConfig{Name: "abc", CustomTemplateFile: ""}, TemplateConfig{Name: "", CustomTemplateFile: "/abc/def"}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := TemplateConfigs{}
			err := yaml.Unmarshal([]byte(tc.input), &got)
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestTemplateConfigFlags(t *testing.T) {
	tests := []struct {
		arguments []string
		expected  TemplateConfigs
	}{
		{
			arguments: []string{"--template", "abc"},
			expected:  TemplateConfigs{TemplateConfig{Name: "abc", CustomTemplateFile: ""}},
		},
		{
			arguments: []string{"-template", "abc"},
			expected:  TemplateConfigs{TemplateConfig{Name: "abc", CustomTemplateFile: ""}},
		},
		{
			arguments: []string{"-template", "abc", "-template", "def"},
			expected:  TemplateConfigs{TemplateConfig{Name: "abc", CustomTemplateFile: ""}, TemplateConfig{Name: "def", CustomTemplateFile: ""}},
		},
		{
			arguments: []string{"-template-file", "abc"},
			expected:  TemplateConfigs{TemplateConfig{Name: "", CustomTemplateFile: "abc"}},
		},
		{
			arguments: []string{"-template-file", "abc", "-template-file", "def"},
			expected:  TemplateConfigs{TemplateConfig{Name: "", CustomTemplateFile: "abc"}, TemplateConfig{Name: "", CustomTemplateFile: "def"}},
		},
		{
			arguments: []string{"-template", "abc", "-template-file", "def"},
			expected:  TemplateConfigs{TemplateConfig{Name: "abc", CustomTemplateFile: ""}, TemplateConfig{Name: "", CustomTemplateFile: "def"}},
		},
		{
			arguments: []string{"-template-file", "abc", "-template", "def"},
			expected:  TemplateConfigs{TemplateConfig{Name: "", CustomTemplateFile: "abc"}, TemplateConfig{Name: "def", CustomTemplateFile: ""}},
		},
		{
			arguments: []string{"--template-file", "abc", "--template", "def"},
			expected:  TemplateConfigs{TemplateConfig{Name: "", CustomTemplateFile: "abc"}, TemplateConfig{Name: "def", CustomTemplateFile: ""}},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fs := flag.NewFlagSet("tc", flag.ContinueOnError)
			c := TemplateConfigs{}
			c.RegisterFlags(fs)
			_ = fs.Parse(tc.arguments)
			assert.Equal(t, tc.expected, c)
		})
	}
}

//go:embed testdata/*
var testData embed.FS

func tdFile(name string) string {
	b, err := testData.ReadFile(path.Join("testdata", name))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func TestConfigEndToEnd(t *testing.T) {
	tests := map[string]struct {
		cliArgs                  []string
		configFileName           string
		configFileContent        string
		defaultConfigFileContent string
		env                      map[string]string
		expected                 Config
		expectedErrorAsserter    tst.ErrorAssertionFunc
	}{
		"default": {
			cliArgs:                  []string{},
			configFileName:           "",
			configFileContent:        "",
			defaultConfigFileContent: "",
			env:                      map[string]string{},
			expected:                 DefaultConfig,
			expectedErrorAsserter:    tst.NoError(),
		},
		"custom config file": {
			cliArgs:                  []string{"--config", "cfg.yml"},
			configFileName:           "cfg.yml",
			configFileContent:        tdFile("cfg1.yml"),
			defaultConfigFileContent: "",
			env:                      map[string]string{},
			expected: Config{
				PackagesQuery: PackagesQueryConfig{
					IncludeTests: false,
					Env:          nil,
					BuildFlags:   nil,
					Dir:          "",
					Patterns:     []string{"./internal/app", "./internal/domain/..."},
				},
				Templates:  TemplateConfigs{TemplateConfig{Name: "otel", CustomTemplateFile: ""}},
				Generate:   GenerateConfig{OutputFile: "zz_generated.{{ .TemplateName }}.go", OutputFileMod: os.FileMode(0o644)},
				Verbose:    false,
				configFile: "cfg.yml",
			},
			expectedErrorAsserter: tst.NoError(),
		},
		"custom config file missing": {
			cliArgs:                  []string{"--config", "cfg.yml"},
			configFileName:           "",
			configFileContent:        "",
			defaultConfigFileContent: "",
			env:                      map[string]string{},
			expected:                 Config{},
			expectedErrorAsserter:    tst.Error(),
		},
		"default config file": {
			cliArgs:                  []string{},
			configFileName:           "",
			configFileContent:        "",
			defaultConfigFileContent: tdFile("cfg1.yml"),
			env:                      map[string]string{},
			expected: Config{
				PackagesQuery: PackagesQueryConfig{
					IncludeTests: false,
					Env:          nil,
					BuildFlags:   nil,
					Dir:          "",
					Patterns:     []string{"./internal/app", "./internal/domain/..."},
				},
				Templates:  TemplateConfigs{TemplateConfig{Name: "otel", CustomTemplateFile: ""}},
				Generate:   GenerateConfig{OutputFile: "zz_generated.{{ .TemplateName }}.go", OutputFileMod: os.FileMode(0o644)},
				Verbose:    false,
				configFile: "",
			},
			expectedErrorAsserter: tst.NoError(),
		},
		"custom config file and cli overwrites": {
			cliArgs:                  []string{"--config", "cfg.yml", "--include_tests", "--template", "pkgpath", "--template-file", "./custom.tmpl"},
			configFileName:           "cfg.yml",
			configFileContent:        tdFile("cfg1.yml"),
			defaultConfigFileContent: "",
			env:                      map[string]string{},
			expected: Config{
				PackagesQuery: PackagesQueryConfig{
					IncludeTests: true,
					Env:          nil,
					BuildFlags:   nil,
					Dir:          "",
					Patterns:     []string{"./internal/app", "./internal/domain/..."},
				},
				Templates:  TemplateConfigs{TemplateConfig{Name: "pkgpath", CustomTemplateFile: ""}, TemplateConfig{Name: "", CustomTemplateFile: "./custom.tmpl"}},
				Generate:   GenerateConfig{OutputFile: "zz_generated.{{ .TemplateName }}.go", OutputFileMod: os.FileMode(0o644)},
				Verbose:    false,
				configFile: "cfg.yml",
			},
			expectedErrorAsserter: tst.NoError(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pwd := t.TempDir()
			t.Chdir(pwd)

			if tc.defaultConfigFileContent != "" {
				fl := path.Join(pwd, ".pkgen.yml")
				err := os.WriteFile(fl, []byte(tc.defaultConfigFileContent), os.FileMode(0o644))
				require.NoError(t, err)
			}

			if tc.configFileName != "" {
				fl := path.Join(pwd, tc.configFileName)
				err := os.WriteFile(fl, []byte(tc.configFileContent), os.FileMode(0o644))
				require.NoError(t, err)
			}

			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			fs := flag.NewFlagSet("flagsettest", flag.PanicOnError)

			got, err := NewConfigGivenCLI(t.Context(), fs, tc.cliArgs)
			tc.expectedErrorAsserter(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
