package config

type Config struct {
	Plugins map[string]map[string]interface{} `yaml:"plugins"`
}
func (c *Config) GetDefaultConfig() string {
    return `
---
plugins:
  git:
    search_path:
      - $HOME/src/*
`
}
