package main

import (
	"HorizonBackend/config"
	"HorizonBackend/internal/handler"
	"HorizonBackend/internal/repository/postgres"
	"HorizonBackend/internal/service"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	conf, confErr := config.Load()
	if confErr != nil {
		fmt.Println("Ошибка загрузки конфигурации:", confErr)
		return
	}

	// Подключение к PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.PgHost, conf.PgPort, conf.PgUser, conf.PgPass, conf.PgDBName, conf.PgSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка при подключении к PostgreSQL:", err)
		return
	}
	defer db.Close()

	imageRepo := postgres.NewImageRepository(db)
	imageService := service.NewImageService(imageRepo)
	// ... (здесь ваш код для инициализации обработчиков и добавления их в роутер)

	r := mux.NewRouter()
	r.HandleFunc("/test", handler.GetImageTest(imageService)).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})

	corsHandler := c.Handler(r)

	fmt.Printf("Сервер запущен на порту %s...\n", conf.Port)
	err = http.ListenAndServe(":"+conf.Port, corsHandler)
	if err != nil {
		return
	}
}
