package main

import (
	"HorizonBackend/config"
	"HorizonBackend/scripts"
	"fmt"
	_ "github.com/lib/pq"
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
	r := NewRouter(db)

	// Запускаем HTTP-сервер на порту 8080
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
