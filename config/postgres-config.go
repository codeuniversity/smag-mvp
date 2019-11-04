package config

import "github.com/codeuniversity/smag-mvp/utils"

// PostgresConfig holds all the configurable variables for Postgres
type PostgresConfig struct {
	PostgresHost     string
	PostgresPassword string
}

//GetPostgresConfig returns a inizialized Postgres Config
func GetPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		PostgresHost:     utils.GetStringFromEnvWithDefault("POSTGRES_HOST", "127.0.0.1"),
		PostgresPassword: utils.GetStringFromEnvWithDefault("POSTGRES_PASSWORD", ""),
	}
}
