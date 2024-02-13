package config

type Config struct {
	AppAddr                string
	AppSelfURL             string
	AdminAddr              string
	AdminSelfURL           string
	Prod                   bool
	PostgresURL            string
	PostgresMigrationsPath string
}

func NewDefault() *Config {
	return &Config{
		AppAddr:                ":9090",
		AppSelfURL:             "http://localhost:9090",
		AdminAddr:              ":9190",
		AdminSelfURL:           "http://localhost:9190",
		Prod:                   false,
		PostgresURL:            "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_pdsdb",
		PostgresMigrationsPath: "file://pds/migrations",
	}
}
