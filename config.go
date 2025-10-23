package pkgen

import (
	"context"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"go.yaml.in/yaml/v4"
)

type RunningMode int

const (
	CLI RunningMode = iota
	GoGenerate
)

func (r RunningMode) String() string {
	switch r {
	case CLI:
		return "cli"
	case GoGenerate:
		return "go-generate"
	default:
		return ""
	}
}

func GetRunningMode() RunningMode {
	_, exists := os.LookupEnv("GOFILE")

	if exists {
		return GoGenerate
	}

	return CLI
}

func runningInsideGoGenerate() bool {
	return GetRunningMode() == GoGenerate
}

func parseYAMLConfig(dst *Config, filePath string) error {
	filePath = filepath.Clean(filePath)

	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if b, err := os.ReadFile(filePath); err == nil {
		err = yaml.Unmarshal(b, &dst)
		if err != nil {
			return err
		}
	}

	return nil
}

var DefaultConfig = Config{
	PackagesQuery: PackagesQueryConfig{
		IncludeTests: false,
		Env:          nil,
		BuildFlags:   nil,
		Dir:          "",
		Patterns:     []string{"./..."},
	},
	Templates: TemplateConfigs{},
	Generate: GenerateConfig{
		OutputFile:    defaultOutputNameTemplate,
		OutputFileMod: os.FileMode(0o644),
	},
	Verbose:    false,
	configFile: "",
}

type Config struct {
	PackagesQuery PackagesQueryConfig `yaml:"packages_query"`
	Templates     TemplateConfigs     `yaml:"templates"`
	Generate      GenerateConfig      `yaml:"generate"`
	Verbose       bool                `yaml:"verbose"`
	configFile    string              // only used when parsing cli arguments
}

func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	c.PackagesQuery.RegisterFlags(fs)
	c.Templates.RegisterFlags(fs)
	c.Generate.RegisterFlags(fs)

	fs.StringVar(&c.configFile, "config", "", "configuration file to use")
	fs.StringVar(&c.configFile, "c", "", "configuration file to use")
	fs.BoolVar(&c.Verbose, "verbose", false, "verbose output")
}

const defaultConfigFile = ".pkgen.yml"

func NewConfig(ctx context.Context) (Config, error) {
	return NewConfigGivenCLI(ctx, flag.CommandLine, os.Args[1:])
}

func NewConfigGivenCLI(ctx context.Context, fs *flag.FlagSet, args []string) (Config, error) {
	cliCnf := Config{}
	cliCnf.RegisterFlags(fs)
	err := fs.Parse(args)
	if err != nil {
		return Config{}, err
	}

	if runningInsideGoGenerate() {
		cnf := merge(cliCnf, DefaultConfig)
		// when running inside go:generate query only the local package. (forced)
		cnf.PackagesQuery.Patterns = []string{"."}
		return cnf, nil
	}

	if cliCnf.configFile != "" {
		fileCnf := Config{}
		err := parseYAMLConfig(&fileCnf, cliCnf.configFile)
		if err != nil {
			return Config{}, err
		}
		return mergeAll(cliCnf, fileCnf, DefaultConfig), nil
	}

	// attempts to read the default but does not fail if the file does not exist.
	defaultFileCnf := Config{}
	err = parseYAMLConfig(&defaultFileCnf, defaultConfigFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	return mergeAll(cliCnf, defaultFileCnf, DefaultConfig), nil
}

type PackagesQueryConfig struct {
	IncludeTests bool     `yaml:"include_tests"`
	Env          []string `yaml:"env"`
	BuildFlags   []string `yaml:"build_flags"`

	Dir      string   `yaml:"dir"`
	Patterns []string `yaml:"patterns"` // e.g. "./..."
}

func (c *PackagesQueryConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.BoolFunc("include_tests", "Include tests in package scanning", func(s string) error {
		c.IncludeTests = true
		return nil
	})
	fs.Func("env", "Additional env var to be passed to package scanning in form of \"key=value\". Can be used multiple times.", func(s string) error {
		c.Env = append(c.Env, s)
		return nil
	})
	fs.Func("build_flag", "Additional build flag to be passed to package scanning. Can be used multiple times.", func(s string) error {
		c.BuildFlags = append(c.BuildFlags, s)
		return nil
	})
	fs.StringVar(&c.Dir, "dir", "", "Dir is the directory in which to run the build system's query tool that provides information about the packages. If Dir is empty, the tool is run in the current directory.")
	fs.Func("pattern", "additional build flags to be passed to package scanning. Can be used multiple times.", func(s string) error {
		c.Patterns = append(c.Patterns, s)
		return nil
	})
}

type GenerateConfig struct {
	OutputFile    string      `yaml:"output"` // the default pattern is zz_generated.{{template name}}.go
	OutputFileMod os.FileMode `yaml:"mod"`
}

func (c *GenerateConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.OutputFile, "output", DefaultConfig.Generate.OutputFile, "The generated file name. The default pattern is zz_generated.{{template name}}.go")
	fs.Func("mod", "The generated file mode in octal format.", func(s string) error {
		oc, err := parseOctal(s)
		if err != nil {
			return err
		}
		c.OutputFileMod = os.FileMode(oc) //nolint:gosec // reason: safe conversion for this use case
		return nil
	})
}

func parseOctal(s string) (uint64, error) {
	if len(s) >= 2 && (s[0:2] == "0o" || s[0:2] == "0O") {
		return strconv.ParseUint(s[2:], 8, 32)
	}
	return strconv.ParseUint(s, 8, 32)
}

type TemplateConfig struct {
	Name               string `yaml:"name"`
	CustomTemplateFile string `yaml:"template_file"`
}

func (tc *TemplateConfig) UnmarshalYAML(value *yaml.Node) error {
	// if it is an object, unmarshal with renaming to avoid cycle.
	if value.Kind == yaml.MappingNode {
		type ObjectTemplateConfig TemplateConfig
		return value.Decode((*ObjectTemplateConfig)(tc))
	}

	// if not an object, it should be a string
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}

	*tc = TemplateConfig{
		Name:               str,
		CustomTemplateFile: "",
	}
	return nil
}

type TemplateConfigs []TemplateConfig

func (tc *TemplateConfigs) RegisterFlags(fs *flag.FlagSet) {
	fs.Func("template", "Add a template to use. Can be used multiple times.", func(s string) error {
		(*tc) = append((*tc), TemplateConfig{Name: s, CustomTemplateFile: ""})
		return nil
	})
	fs.Func("template-file", "Add a path to a custom template to use. Can be used multiple times.", func(s string) error {
		(*tc) = append((*tc), TemplateConfig{Name: "", CustomTemplateFile: s})
		return nil
	})
}

func (tc *TemplateConfigs) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		type SequenceTemplateConfigs TemplateConfigs
		return value.Decode((*SequenceTemplateConfigs)(tc))
	}

	// if it is a single item, decode as such
	c := TemplateConfig{}
	err := value.Decode(&c)
	if err != nil {
		return err
	}

	*tc = []TemplateConfig{c}
	return nil
}

var (
	_ yaml.Unmarshaler = (*TemplateConfigs)(nil)
	_ yaml.Unmarshaler = (*TemplateConfig)(nil)
)

var ErrMalformedTemplateConfigs = errors.New("malformed template configs")

//
// merge
//

func mergeAll(c ...Config) Config {
	if len(c) == 0 {
		return Config{}
	}

	cfg := Config{}

	for _, cn := range c {
		cfg = merge(cfg, cn)
	}

	return cfg
}

func merge(a, b Config) Config {
	return Config{
		PackagesQuery: PackagesQueryConfig{
			IncludeTests: firstNotEmpty(a.PackagesQuery.IncludeTests, b.PackagesQuery.IncludeTests),
			Env:          firstNotEmptySlice(a.PackagesQuery.Env, b.PackagesQuery.Env),
			BuildFlags:   firstNotEmptySlice(a.PackagesQuery.BuildFlags, b.PackagesQuery.BuildFlags),
			Dir:          firstNotEmpty(a.PackagesQuery.Dir, b.PackagesQuery.Dir),
			Patterns:     firstNotEmptySlice(a.PackagesQuery.Patterns, b.PackagesQuery.Patterns),
		},
		Templates: firstNotEmptySlice(a.Templates, b.Templates),
		Generate: GenerateConfig{
			OutputFile:    firstNotEmpty(a.Generate.OutputFile, b.Generate.OutputFile),
			OutputFileMod: firstNotEmpty(a.Generate.OutputFileMod, b.Generate.OutputFileMod),
		},
		Verbose:    firstNotEmpty(a.Verbose, b.Verbose),
		configFile: firstNotEmpty(a.configFile, b.configFile),
	}
}

func firstNotEmpty[T comparable](a, b T) T { //nolint: ireturn
	var empty T
	if a == empty {
		return b
	}

	return a
}

func firstNotEmptySlice[S ~[]E, E any](a, b S) S { //nolint: ireturn
	if len(a) > 0 {
		return a
	}

	return b
}
