package go_rest_starter_pack

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Configuration struct {
		Server   *Server
		Database *Database
		JWT      *JWT
	}

	Server struct {
		Port string
	}

	Database struct {
		URI     string
		Name    string
		Timeout int
	}

	JWT struct {
		SecretKey          string
		SigningAlgorithm   string
		AccessTokenExpiry  int
		RefreshTokenExpiry int
	}
)

func Load() (*Configuration, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("configuration error : " + err.Error())
	}

	databaseURI := "mongodb://" + getEnv("MONGODB_USERNAME", "")
	databaseURI += ":" + getEnv("MONGODB_PASSWORD", "")
	databaseURI += "@" + getEnv("MONGODB_HOST", "localhost")
	databaseURI += ":" + getEnv("MONGODB_PORT", "27017")

	timeout, _ := strconv.Atoi(getEnv("DB_TIMEOUT", "5"))
	accessTokenExpiry, _ := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "1"))
	refreshTokenExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "1"))

	return &Configuration{
		Server: &Server{
			Port: ":" + getEnv("SERVER_PORT", "8080"),
		},
		Database: &Database{
			URI:     databaseURI,
			Name:    getEnv("MONGODB_DB_NAME", ""),
			Timeout: timeout,
		},
		JWT: &JWT{
			SecretKey:          getEnv("JWT_SECRET_KEY", ""),
			SigningAlgorithm:   getEnv("JWT_SIGNING_ALGORITHM", ""),
			AccessTokenExpiry:  accessTokenExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
		},
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
