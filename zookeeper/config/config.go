package config

type Config struct {
	AppAddr                string
	AdminAddr              string
	Prod                   bool
	PostgresURL            string
	PostgresMigrationsPath string
	JWTSecret              []byte
}

func NewDefault() *Config {
	return &Config{
		AppAddr:                ":9092",
		AdminAddr:              ":9192",
		Prod:                   false,
		PostgresURL:            "postgresql://ecumenosuser:rootpassword@localhost:5432/ecumenos_zookeeperdb",
		PostgresMigrationsPath: "file://zookeeper/migrations",
		JWTSecret:              []byte("jwtsecretplaceholder"),
	}
}
