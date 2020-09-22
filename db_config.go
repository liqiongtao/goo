package goo

type DBConfig struct {
	Driver   string   `yaml:"driver"`
	Master   string   `yaml:"master"`
	Slaves   []string `yaml:"slaves"`
	LogModel bool     `yaml:"log_model"`
	MaxIdle  int      `yaml:"max_idle"`
	MaxOpen  int      `yaml:"max_open"`
}
