package store

import "os"

type Config struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func ConfigFromEnv() Config {
	return Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Name: os.Getenv("DB_NAME"),
	}
}
