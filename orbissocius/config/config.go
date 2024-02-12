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
		AppAddr:                ":9091",
		AppSelfURL:             "http://localhost:9091",
		AdminAddr:              ":9191",
		AdminSelfURL:           "http://localhost:9191",
		Prod:                   false,
		PostgresURL:            "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_orbissociusdb",
		PostgresMigrationsPath: "file://orbissocius/migrations",
	}
}
