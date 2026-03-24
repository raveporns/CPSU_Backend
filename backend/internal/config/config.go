package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort          string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string

	MinioEndpoint      string
	MinioAccessKey     string
	MinioSecretKey     string
	MinioBucket        string
	MinioUseSSL        bool
	MinioPublicBaseURL string

	CalendarID string
}

func LoadConfig() (Config, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	viper.SetDefault("POSTGRES.HOST", "localhost")
	viper.SetDefault("POSTGRES.PORT", 5432)
	viper.SetDefault("POSTGRES.USER", "postgres")
	viper.SetDefault("POSTGRES.PASSWORD", "")
	viper.SetDefault("POSTGRES.DBNAME", "cpsu")
	viper.SetDefault("POSTGRES.SSLMODE", "disable")

	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin123")
	viper.SetDefault("MINIO_BUCKET", "images")
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("MINIO_PUBLIC_BASE_URL", "http://localhost:9000")

	viper.SetDefault("CALENDAR.ID", "")

	useSSL := viper.GetBool("MINIO_USE_SSL")

	// Set config values
	config := Config{
		AppPort:            viper.GetString("APP.PORT"),
		DatabaseHost:       viper.GetString("POSTGRES.HOST"),
		DatabasePort:       viper.GetInt("POSTGRES.PORT"),
		DatabaseUser:       viper.GetString("POSTGRES.USER"),
		DatabasePassword:   viper.GetString("POSTGRES.PASSWORD"),
		DatabaseName:       viper.GetString("POSTGRES.DBNAME"),
		DatabaseSSLMode:    viper.GetString("POSTGRES.SSLMODE"),
		MinioEndpoint:      viper.GetString("MINIO_ENDPOINT"),
		MinioAccessKey:     viper.GetString("MINIO_ACCESS_KEY"),
		MinioSecretKey:     viper.GetString("MINIO_SECRET_KEY"),
		MinioBucket:        viper.GetString("MINIO_BUCKET"),
		MinioUseSSL:        useSSL,
		MinioPublicBaseURL: viper.GetString("MINIO_PUBLIC_BASE_URL"),
		CalendarID:         viper.GetString("CALENDAR.ID"),
	}

	return config, nil
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode)
}
