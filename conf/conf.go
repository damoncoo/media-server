package conf

// Config 配置
type Config struct {
	Path []string `yaml:"path"`
}

var (
	Conf Config
)
