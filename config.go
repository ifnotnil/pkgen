package pkgen

type Config struct {
	PackagesQuery PackagesQueryConfig `yaml:"packages_query"`
	Templates     string              `yaml:"packages_query"`
}
