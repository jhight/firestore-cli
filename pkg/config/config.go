package config

type Config struct {
	ServiceAccount string       `yaml:"service-account"`
	ProjectID      string       `yaml:"project-id"`
	PrettyPrint    bool         `yaml:"pretty-print"`
	RawPrint       bool         `yaml:"raw"`
	PrettySpacing  int          `yaml:"spacing"`
	Backup         BackupConfig `yaml:"backup"`
	Flatten        bool         `yaml:"flatten"`
}

type BackupConfig struct {
	Collection string   `yaml:"collection"`
	Commands   []string `yaml:"commands"`
}
