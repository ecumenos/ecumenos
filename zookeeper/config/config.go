package config

type Config struct {
	AppAddr                string
	AppSelfURL             string
	AdminAddr              string
	AdminSelfURL           string
	Prod                   bool
	PostgresURL            string
	PostgresMigrationsPath string
	JWTSecret              []byte
}

func NewDefault() *Config {
	return &Config{
		AppAddr:                ":9092",
		AppSelfURL:             "http://localhost:9092",
		AdminAddr:              ":9192",
		AdminSelfURL:           "http://localhost:9192",
		Prod:                   false,
		PostgresURL:            "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_zookeeperdb",
		PostgresMigrationsPath: "file://zookeeper/migrations",
		JWTSecret:              []byte("jwtsecretplaceholder"),
	}
}
