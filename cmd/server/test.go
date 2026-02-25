package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}
	return nil
}

func checkEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", key)
	}
	return value, nil
}

func init() {
	loadEnv()
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
