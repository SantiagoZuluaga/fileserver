package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	TCP_HOST = "localhost"
	TCP_PORT = "5000"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	if os.Getenv("TCP_HOST") != "" {
		TCP_HOST = os.Getenv("TCP_HOST")
	}
	if os.Getenv("TCP_PORT") != "" {
		TCP_PORT = os.Getenv("TCP_PORT")
	}
}
