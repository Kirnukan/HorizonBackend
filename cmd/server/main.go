package main

import (
	"HorizonBackend/config"
	"HorizonBackend/internal/router"
	"HorizonBackend/scripts"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := config.NewConnection(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing the database:", err)
		}
	}()

	// Этот скрипт добавляет изображения из папки в вашу базу данных
	scripts.AddImagesFromFolder(db, "./static/images")

	// Создаем роутер
	r := router.NewRouter(db, cfg)

	// Запускаем HTTP-сервер на порту 8080
	fmt.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
