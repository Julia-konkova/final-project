package server

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Пытаемся загрузить из .env файла, но не падаем если его нет
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found, using system environment variables")
	}
}

func GetValidatedPort() (string, error) {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		portStr = "7540"
		log.Printf("TODO_PORT not set, using default: %s", portStr)
	}

	// Валидация порта
	_, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("TODO_PORT must be a number, got: %s", portStr)
	}

	return portStr, nil
}
