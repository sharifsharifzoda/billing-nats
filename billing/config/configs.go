package config

type NatsConnConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConnConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}
