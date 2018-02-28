package config

var (
	Conf *Config
)

const (
	DBName     = "data"
	TestDBName = "testdata"
)

type Config struct {
	DataDir string
	TestNet bool
}
