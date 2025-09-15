package main

import (
	"final-project/handlers"
	"final-project/pkg/db"
	"final-project/server"
	"log"
)

func main() {
	webDir := "./web"

	addr, err := db.GetDBAddress()
	if err != nil {
		log.Fatalf("Failed to receive database address: %v", err)
		addr = "scheduler.db"
	}
	err = db.Init(addr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DB.Close() // Закрываем соединение при завершении

	// Создаем экземпляр хендлеров
	h := handlers.New(webDir)

	// Создаем и настраиваем сервер
	server.LoadEnv()

	// Получаем и валидируем порт
	port, err := server.GetValidatedPort()
	if err != nil {
		log.Fatalf("Invalid port configuration: %v", err)
	}

	srv := server.New(":" + port)
	srv.SetupRoutes(h)

	// Запускаем сервер
	log.Printf("Server starting on port %s", port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
