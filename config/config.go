package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port string
}

var config = Config{
	Host: "localhost",
	Port: "5000",
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	if os.Getenv("TCP_HOST") != "" {
		config.Host = os.Getenv("TCP_HOST")
	}
	if os.Getenv("TCP_PORT") != "" {
		config.Port = os.Getenv("TCP_PORT")
	}
}

func GetConfig() Config {
	return config
}
