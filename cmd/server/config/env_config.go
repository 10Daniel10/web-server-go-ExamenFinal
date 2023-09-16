package config

import (
	"fmt"
	"os"
)

var envs = map[string]PublicConfig{
	"local": {
		PubKey: "local_key",
	},
	"dev": {
		PubKey: "dev_key",
	},
	"prod": {
		PubKey: "prod_key",
	},
}

type EnvConfig struct {
	Public  PublicConfig
	Private PrivateConfig
}

type PublicConfig struct {
	PubKey string
}

type PrivateConfig struct {
	// Web server config
	SecretKey string
	Address   string
	Port      string
	Host      string
	BasePath  string
	// DB config
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

func NewEnvConfig(env string) (*EnvConfig, error) {
	publicConfig, ok := envs[env]
	if ok == false {
		return nil, fmt.Errorf("env not found")
	}

	// Private config, web server
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY not found")
	}

	address := os.Getenv("ADDRESS")
	if address == "" {
		return nil, fmt.Errorf("ADDRESS not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT not found")
	}

	host := os.Getenv("HOST")
	if host == "" {
		return nil, fmt.Errorf("HOST not found")
	}

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		return nil, fmt.Errorf("BASE_PATH not found")
	}

	// Private config, database
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("DB_USER not found")
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		return nil, fmt.Errorf("DB_PASS not found")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST not found")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("DB_PORT not found")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME not found")
	}

	return &EnvConfig{
		Public: publicConfig,
		Private: PrivateConfig{
			// Web server config
			SecretKey: secretKey,
			Address:   address,
			Port:      port,
			Host:      host,
			BasePath:  basePath,

			// DB config
			DBUser: dbUser,
			DBPass: dbPass,
			DBHost: dbHost,
			DBPort: dbPort,
			DBName: dbName,
		},
	}, nil
}
