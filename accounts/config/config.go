package config

type Config struct {
	AppAddr    string
	AppSelfURL string
	Prod       bool
}

func NewDefault() *Config {
	return &Config{
		AppAddr:    ":9093",
		AppSelfURL: "http://localhost:9093",
		Prod:       false,
	}
}
