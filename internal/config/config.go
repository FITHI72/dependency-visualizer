package config

import "flag"

type Config struct {
	Name     string
	Repo     string
	TestMode string
	Version  string
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.Name, "name", "", "Имя анализируемого пакета (обязательно)")
	flag.StringVar(&cfg.Repo, "repo", "", "URL или путь к репозиторию (обязательно)")
	flag.StringVar(&cfg.TestMode, "test-mode", "off", "Режим работы: off|local|remote")
	flag.StringVar(&cfg.Version, "version", "", "Версия пакета (например 1.2.3 или latest)")
	flag.Parse()
	return cfg, nil
}
