package configs

type (
	Config struct {
		Service  Service
		Database DatabaseConfig
	}

	Service struct {
		Port      string
		SecretKey string
	}

	DatabaseConfig struct {
		DataSourceName string
	}
)
