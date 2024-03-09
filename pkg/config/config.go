package config

type Config struct {
	ServiceAccount string       `yaml:"service-account"`
	ProjectID      string       `yaml:"project-id"`
	PrettyPrint    bool         `yaml:"pretty-print"`
	Raw            bool         `yaml:"raw"`
	PrettySpacing  int          `yaml:"spacing"`
	Backup         BackupConfig `yaml:"backup"`
}

type BackupConfig struct {
	Collection string   `yaml:"collection"`
	Commands   []string `yaml:"commands"`
}
