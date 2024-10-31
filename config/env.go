package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	Email                  string
	EmailPassword          string
	Hostsmtp               string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	AesKey                 string
	DBName                 string
	JWTSecret              string
	GoogleClientID         string
	GoogleClientSecretId   string
	JWTExpirationInSeconds int64
	TokenExpirationInHour  uint8
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Email:                  getEnv("EMAIL", "EXAMPLE@EMAIL"),
		EmailPassword:          getEnv("EMAIL_PASSWORD", "your super secret code"),
		Hostsmtp:               getEnv("HOST_SMTP", "smtp.mail.com"),
		JWTSecret:              getEnv("JWT_secret", "your super secret code"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_Expiration_In_Seconds", 3600*24*7),
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", " db password"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", " ip host"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "seed"),
		GoogleClientID:         getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecretId:   getEnv("GOOGLE_CLIENT_SECRET_ID", ""),
		AesKey:                 getEnv("AES_KEY", "your super secret code"),
		TokenExpirationInHour:  uint8(getEnvAsInt("TOKEN_EXPIRATION_HOUR", 5)),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
