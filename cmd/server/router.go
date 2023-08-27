package main

import (
	"HorizonBackend/internal/handler"
	"HorizonBackend/internal/repository/postgres"
	"HorizonBackend/internal/service"
	"database/sql"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Инициализация репозитория и службы
	imageRepo := postgres.NewImageRepository(db)
	imageService := service.NewImageService(imageRepo)

	// Регистрация обработчика
	r.HandleFunc("/{family}/{group}/", handler.GetImagesByFamilyAndGroup(imageService)).Methods("GET")
	r.HandleFunc("/image/{id:[0-9]+}", handler.GetImageByID(imageService)).Methods("GET")
	//можете использовать этот роут для поиска изображений,
	//передавая в запросе параметры keyword и family.
	//Например: /search?keyword=Form&family=Forms.
	r.HandleFunc("/search", handler.SearchImages(imageService)).Methods("GET")

	return r
}
