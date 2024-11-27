package config

import (
	"fmt"
	"os"

	"song-library/internal/constants"
)

type Config struct {
	DB            DatabaseConfig
	ServerAddress string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	dbConfig := DatabaseConfig{}

	var err error
	if dbConfig.Host, err = getRequiredEnv(constants.EnvDBHost); err != nil {
		return nil, err
	}
	if dbConfig.Port, err = getRequiredEnv(constants.EnvDBPort); err != nil {
		return nil, err
	}
	if dbConfig.User, err = getRequiredEnv(constants.EnvDBUser); err != nil {
		return nil, err
	}
	if dbConfig.Password, err = getRequiredEnv(constants.EnvDBPassword); err != nil {
		return nil, err
	}
	if dbConfig.DBName, err = getRequiredEnv(constants.EnvDBName); err != nil {
		return nil, err
	}
	if dbConfig.SSLMode, err = getRequiredEnv(constants.EnvDBSSLMode); err != nil {
		return nil, err
	}

	serverHost, err := getRequiredEnv(constants.EnvServerHost)
	if err != nil {
		return nil, err
	}
	serverPort, err := getRequiredEnv(constants.EnvServerPort)
	if err != nil {
		return nil, err
	}
	serverProtocol, err := getRequiredEnv(constants.EnvServerProtocol)
	if err != nil {
		return nil, err
	}
	serverAddress := fmt.Sprintf(constants.DefaultAddressFormat, serverProtocol, serverHost, serverPort)

	return &Config{
		DB:            dbConfig,
		ServerAddress: serverAddress,
	}, nil
}

// возвращает ошибку, если переменная не установлена
func getRequiredEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf(constants.ErrMissingEnvVar, key)
	}
	return value, nil
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf(
		constants.PostgresConnectionString,
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.DBName,
		c.DB.SSLMode,
	)
}
