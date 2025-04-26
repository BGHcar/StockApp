package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	URL       string
	User      string
	Pass      string
	Host      string
	Port      string
	Name      string
	Mode      string
	TableName string
}

type ExternalApi struct {
	URL   string
	Token string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Theres no environment")
	}
}

func LoadDB() DBConfig {
	return DBConfig{
		URL:       os.Getenv("DATABASE_URL"),
		User:      os.Getenv("SQL_USER"),
		Pass:      os.Getenv("GENERATED_PASSWORD"),
		Host:      os.Getenv("CLUSTER_HOST"),
		Port:      os.Getenv("CLUSTER_PORT"),
		Name:      os.Getenv("CLUSTER_NAME"),
		Mode:      os.Getenv("DB_SSL_MODE"),
		TableName: os.Getenv("TABLE_NAME"),
	}
}

func LoadApi() ExternalApi {
	return ExternalApi{
		URL:   os.Getenv("API_URL"),
		Token: os.Getenv("API_TOKEN"),
	}
}

func LoadPort() string {
	return os.Getenv("PORT")
}
